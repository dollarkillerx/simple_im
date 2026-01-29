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

// ============ message.send ============

type MessageSendMethod struct {
	storage *storage.Storage
	hub     *ws.Hub
}

func NewMessageSendMethod(s *storage.Storage, h *ws.Hub) *MessageSendMethod {
	return &MessageSendMethod{storage: s, hub: h}
}

func (m *MessageSendMethod) Name() string { return "message.send" }

func (m *MessageSendMethod) RequireAuth() bool { return true }

type MessageSendParams struct {
	ReceiverID int64             `json:"receiver_id"` // For private chat
	GroupID    int64             `json:"group_id"`    // For group chat
	MsgType    models.MessageType `json:"msg_type"`    // 1:text 2:image 3:file
	Content    string            `json:"content"`
	FileURL    string            `json:"file_url"`
	FileName   string            `json:"file_name"`
	FileSize   int64             `json:"file_size"`
}

func (m *MessageSendMethod) Execute(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var p MessageSendParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, fmt.Errorf("invalid params: %v", err)
	}

	if p.ReceiverID == 0 && p.GroupID == 0 {
		return nil, errors.New("receiver_id or group_id is required")
	}

	if p.MsgType == 0 {
		p.MsgType = models.MsgTypeText
	}

	if p.MsgType == models.MsgTypeText && p.Content == "" {
		return nil, errors.New("content is required for text message")
	}

	if (p.MsgType == models.MsgTypeImage || p.MsgType == models.MsgTypeFile) && p.FileURL == "" {
		return nil, errors.New("file_url is required for image/file message")
	}

	userID := ctx.Value("user_id").(int64)
	username := ctx.Value("username").(string)
	db := m.storage.GetDB()

	// Validate receiver or group
	var groupMembers []int64
	if p.GroupID > 0 {
		// Check if user is member of group
		var membership models.GroupMember
		err := db.Where("group_id = ? AND user_id = ?", p.GroupID, userID).First(&membership).Error
		if err != nil {
			return nil, errors.New("not a member of this group")
		}

		// Get all group members for broadcasting
		var members []models.GroupMember
		db.Where("group_id = ?", p.GroupID).Find(&members)
		for _, m := range members {
			groupMembers = append(groupMembers, m.UserID)
		}
	} else {
		// Check if receiver exists and is friend
		var friend models.Friend
		err := db.Where("((user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)) AND status = ?",
			userID, p.ReceiverID, p.ReceiverID, userID, models.FriendStatusAccepted).First(&friend).Error
		if err != nil {
			return nil, errors.New("can only send messages to friends")
		}
	}

	// Create message
	msg := &models.Message{
		SenderID:   userID,
		ReceiverID: p.ReceiverID,
		GroupID:    p.GroupID,
		MsgType:    p.MsgType,
		Content:    p.Content,
		FileURL:    p.FileURL,
		FileName:   p.FileName,
		FileSize:   p.FileSize,
		CreatedAt:  time.Now(),
	}

	if err := db.Create(msg).Error; err != nil {
		return nil, fmt.Errorf("failed to create message: %v", err)
	}

	// Get group name if group chat
	var groupName string
	if p.GroupID > 0 {
		var group models.Group
		db.First(&group, p.GroupID)
		groupName = group.Name
	}

	// Broadcast message via WebSocket
	m.hub.Broadcast(&ws.Message{
		ID:           msg.ID,
		Type:         "message",
		SenderID:     userID,
		SenderName:   username,
		ReceiverID:   p.ReceiverID,
		GroupID:      p.GroupID,
		GroupName:    groupName,
		MsgType:      ws.MessageType(p.MsgType),
		Content:      p.Content,
		FileURL:      p.FileURL,
		FileName:     p.FileName,
		FileSize:     p.FileSize,
		CreatedAt:    msg.CreatedAt,
		GroupMembers: groupMembers,
	})

	return msg, nil
}

// ============ message.history ============

type MessageHistoryMethod struct {
	storage *storage.Storage
}

func NewMessageHistoryMethod(s *storage.Storage) *MessageHistoryMethod {
	return &MessageHistoryMethod{storage: s}
}

func (m *MessageHistoryMethod) Name() string { return "message.history" }

func (m *MessageHistoryMethod) RequireAuth() bool { return true }

type MessageHistoryParams struct {
	ReceiverID int64 `json:"receiver_id"` // For private chat
	GroupID    int64 `json:"group_id"`    // For group chat
	BeforeID   int64 `json:"before_id"`   // For pagination
	Limit      int   `json:"limit"`
}

func (m *MessageHistoryMethod) Execute(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var p MessageHistoryParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, fmt.Errorf("invalid params: %v", err)
	}

	if p.ReceiverID == 0 && p.GroupID == 0 {
		return nil, errors.New("receiver_id or group_id is required")
	}

	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 50
	}

	userID := ctx.Value("user_id").(int64)
	db := m.storage.GetDB()

	var messages []models.Message
	query := db.Preload("Sender").Order("id DESC").Limit(p.Limit)

	if p.BeforeID > 0 {
		query = query.Where("id < ?", p.BeforeID)
	}

	if p.GroupID > 0 {
		// Check if user is member of group
		var membership models.GroupMember
		err := db.Where("group_id = ? AND user_id = ?", p.GroupID, userID).First(&membership).Error
		if err != nil {
			return nil, errors.New("not a member of this group")
		}

		query = query.Where("group_id = ?", p.GroupID)
	} else {
		// Private chat: messages between two users
		query = query.Where(
			"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID, p.ReceiverID, p.ReceiverID, userID,
		)
	}

	if err := query.Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("failed to get messages: %v", err)
	}

	// Reverse to chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
