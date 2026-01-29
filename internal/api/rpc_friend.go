package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"simple_im/internal/models"
	"simple_im/internal/storage"
	"simple_im/internal/ws"
)

// ============ friend.list ============

type FriendListMethod struct {
	storage *storage.Storage
}

func NewFriendListMethod(s *storage.Storage) *FriendListMethod {
	return &FriendListMethod{storage: s}
}

func (m *FriendListMethod) Name() string { return "friend.list" }

func (m *FriendListMethod) RequireAuth() bool { return true }

func (m *FriendListMethod) Execute(ctx context.Context, params json.RawMessage) (interface{}, error) {
	userID := ctx.Value("user_id").(int64)
	db := m.storage.GetDB()

	var friends []models.Friend
	err := db.Where("(user_id = ? OR friend_id = ?) AND status = ?", userID, userID, models.FriendStatusAccepted).
		Preload("User").
		Preload("Friend").
		Find(&friends).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get friends: %v", err)
	}

	// Build friend list with user info
	result := make([]map[string]interface{}, 0, len(friends))
	for _, f := range friends {
		var friendUser *models.User
		if f.UserID == userID {
			friendUser = f.Friend
		} else {
			friendUser = f.User
		}
		if friendUser != nil {
			result = append(result, map[string]interface{}{
				"id":         f.ID,
				"user_id":    friendUser.ID,
				"username":   friendUser.Username,
				"nickname":   friendUser.Nickname,
				"avatar":     friendUser.Avatar,
				"created_at": f.CreatedAt,
			})
		}
	}

	return result, nil
}

// ============ friend.add ============

type FriendAddMethod struct {
	storage *storage.Storage
	hub     *ws.Hub
}

func NewFriendAddMethod(s *storage.Storage, h *ws.Hub) *FriendAddMethod {
	return &FriendAddMethod{storage: s, hub: h}
}

func (m *FriendAddMethod) Name() string { return "friend.add" }

func (m *FriendAddMethod) RequireAuth() bool { return true }

type FriendAddParams struct {
	FriendID int64  `json:"friend_id"`
	Username string `json:"username"`
}

func (m *FriendAddMethod) Execute(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var p FriendAddParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, fmt.Errorf("invalid params: %v", err)
	}

	userID := ctx.Value("user_id").(int64)
	username := ctx.Value("username").(string)
	db := m.storage.GetDB()

	friendID := p.FriendID

	// If username provided, find user by username
	if friendID == 0 && p.Username != "" {
		var user models.User
		if err := db.Where("username = ?", p.Username).First(&user).Error; err != nil {
			return nil, errors.New("user not found")
		}
		friendID = user.ID
	}

	if friendID == 0 {
		return nil, errors.New("friend_id or username is required")
	}

	if friendID == userID {
		return nil, errors.New("cannot add yourself as friend")
	}

	// Check if already friends or pending
	var existing models.Friend
	err := db.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
		userID, friendID, friendID, userID).First(&existing).Error
	if err == nil {
		if existing.Status == models.FriendStatusAccepted {
			return nil, errors.New("already friends")
		}
		if existing.Status == models.FriendStatusPending {
			return nil, errors.New("friend request already pending")
		}
	}

	friend := &models.Friend{
		UserID:   userID,
		FriendID: friendID,
		Status:   models.FriendStatusPending,
	}

	if err := db.Create(friend).Error; err != nil {
		return nil, fmt.Errorf("failed to create friend request: %v", err)
	}

	// Send notification via WebSocket
	m.hub.Broadcast(&ws.Message{
		Type:       "friend_request",
		SenderID:   userID,
		SenderName: username,
		ReceiverID: friendID,
		Content:    "sent you a friend request",
		CreatedAt:  time.Now(),
	})

	return map[string]interface{}{
		"id":      friend.ID,
		"message": "friend request sent",
	}, nil
}

// ============ friend.accept ============

type FriendAcceptMethod struct {
	storage *storage.Storage
	hub     *ws.Hub
}

func NewFriendAcceptMethod(s *storage.Storage, h *ws.Hub) *FriendAcceptMethod {
	return &FriendAcceptMethod{storage: s, hub: h}
}

func (m *FriendAcceptMethod) Name() string { return "friend.accept" }

func (m *FriendAcceptMethod) RequireAuth() bool { return true }

type FriendAcceptParams struct {
	RequestID int64 `json:"request_id"`
	Accept    bool  `json:"accept"`
}

func (m *FriendAcceptMethod) Execute(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var p FriendAcceptParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, fmt.Errorf("invalid params: %v", err)
	}

	if p.RequestID == 0 {
		return nil, errors.New("request_id is required")
	}

	userID := ctx.Value("user_id").(int64)
	username := ctx.Value("username").(string)
	db := m.storage.GetDB()

	var friend models.Friend
	if err := db.First(&friend, p.RequestID).Error; err != nil {
		return nil, errors.New("friend request not found")
	}

	if friend.FriendID != userID {
		return nil, errors.New("not authorized to accept this request")
	}

	if friend.Status != models.FriendStatusPending {
		return nil, errors.New("request already processed")
	}

	if p.Accept {
		friend.Status = models.FriendStatusAccepted
	} else {
		friend.Status = models.FriendStatusRejected
	}

	if err := db.Save(&friend).Error; err != nil {
		return nil, fmt.Errorf("failed to update friend request: %v", err)
	}

	// Notify the requester
	notifyType := "friend_accepted"
	content := "accepted your friend request"
	if !p.Accept {
		notifyType = "friend_rejected"
		content = "rejected your friend request"
	}

	m.hub.Broadcast(&ws.Message{
		Type:       notifyType,
		SenderID:   userID,
		SenderName: username,
		ReceiverID: friend.UserID,
		Content:    content,
		CreatedAt:  time.Now(),
	})

	return map[string]interface{}{
		"message": "friend request " + notifyType[7:],
	}, nil
}

// ============ friend.pending ============

type FriendPendingMethod struct {
	storage *storage.Storage
}

func NewFriendPendingMethod(s *storage.Storage) *FriendPendingMethod {
	return &FriendPendingMethod{storage: s}
}

func (m *FriendPendingMethod) Name() string { return "friend.pending" }

func (m *FriendPendingMethod) RequireAuth() bool { return true }

func (m *FriendPendingMethod) Execute(ctx context.Context, params json.RawMessage) (interface{}, error) {
	userID := ctx.Value("user_id").(int64)
	db := m.storage.GetDB()

	var requests []models.Friend
	err := db.Where("friend_id = ? AND status = ?", userID, models.FriendStatusPending).
		Preload("User").
		Find(&requests).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get pending requests: %v", err)
	}

	result := make([]map[string]interface{}, 0, len(requests))
	for _, r := range requests {
		result = append(result, map[string]interface{}{
			"id":         r.ID,
			"user_id":    r.UserID,
			"username":   r.User.Username,
			"nickname":   r.User.Nickname,
			"avatar":     r.User.Avatar,
			"created_at": r.CreatedAt,
		})
	}

	return result, nil
}
