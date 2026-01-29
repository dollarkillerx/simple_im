package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"simple_im/internal/models"
	"simple_im/internal/storage"
)

// ============ group.create ============

type GroupCreateMethod struct {
	storage *storage.Storage
}

func NewGroupCreateMethod(s *storage.Storage) *GroupCreateMethod {
	return &GroupCreateMethod{storage: s}
}

func (m *GroupCreateMethod) Name() string { return "group.create" }

func (m *GroupCreateMethod) RequireAuth() bool { return true }

type GroupCreateParams struct {
	Name      string  `json:"name"`
	Avatar    string  `json:"avatar"`
	MemberIDs []int64 `json:"member_ids"`
}

func (m *GroupCreateMethod) Execute(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var p GroupCreateParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, fmt.Errorf("invalid params: %v", err)
	}

	if p.Name == "" {
		return nil, errors.New("group name is required")
	}

	userID := ctx.Value("user_id").(int64)
	db := m.storage.GetDB()

	// Create group
	group := &models.Group{
		Name:    p.Name,
		OwnerID: userID,
		Avatar:  p.Avatar,
	}

	tx := db.Begin()

	if err := tx.Create(group).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create group: %v", err)
	}

	// Add owner as member
	ownerMember := &models.GroupMember{
		GroupID:  group.ID,
		UserID:   userID,
		Role:     models.GroupRoleOwner,
		JoinedAt: time.Now(),
	}
	if err := tx.Create(ownerMember).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to add owner as member: %v", err)
	}

	// Add other members
	for _, memberID := range p.MemberIDs {
		if memberID == userID {
			continue
		}
		member := &models.GroupMember{
			GroupID:  group.ID,
			UserID:   memberID,
			Role:     models.GroupRoleMember,
			JoinedAt: time.Now(),
		}
		if err := tx.Create(member).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to add member: %v", err)
		}
	}

	tx.Commit()

	return group, nil
}

// ============ group.list ============

type GroupListMethod struct {
	storage *storage.Storage
}

func NewGroupListMethod(s *storage.Storage) *GroupListMethod {
	return &GroupListMethod{storage: s}
}

func (m *GroupListMethod) Name() string { return "group.list" }

func (m *GroupListMethod) RequireAuth() bool { return true }

func (m *GroupListMethod) Execute(ctx context.Context, params json.RawMessage) (interface{}, error) {
	userID := ctx.Value("user_id").(int64)
	db := m.storage.GetDB()

	var members []models.GroupMember
	err := db.Where("user_id = ?", userID).
		Preload("Group").
		Preload("Group.Owner").
		Find(&members).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get groups: %v", err)
	}

	result := make([]map[string]interface{}, 0, len(members))
	for _, m := range members {
		if m.Group != nil {
			result = append(result, map[string]interface{}{
				"id":         m.Group.ID,
				"name":       m.Group.Name,
				"avatar":     m.Group.Avatar,
				"owner_id":   m.Group.OwnerID,
				"owner_name": m.Group.Owner.Nickname,
				"role":       m.Role,
				"joined_at":  m.JoinedAt,
			})
		}
	}

	return result, nil
}

// ============ group.info ============

type GroupInfoMethod struct {
	storage *storage.Storage
}

func NewGroupInfoMethod(s *storage.Storage) *GroupInfoMethod {
	return &GroupInfoMethod{storage: s}
}

func (m *GroupInfoMethod) Name() string { return "group.info" }

func (m *GroupInfoMethod) RequireAuth() bool { return true }

type GroupInfoParams struct {
	GroupID int64 `json:"group_id"`
}

func (m *GroupInfoMethod) Execute(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var p GroupInfoParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, fmt.Errorf("invalid params: %v", err)
	}

	if p.GroupID == 0 {
		return nil, errors.New("group_id is required")
	}

	userID := ctx.Value("user_id").(int64)
	db := m.storage.GetDB()

	// Check if user is member of group
	var membership models.GroupMember
	err := db.Where("group_id = ? AND user_id = ?", p.GroupID, userID).First(&membership).Error
	if err != nil {
		return nil, errors.New("not a member of this group")
	}

	var group models.Group
	err = db.Preload("Owner").Preload("Members").Preload("Members.User").First(&group, p.GroupID).Error
	if err != nil {
		return nil, errors.New("group not found")
	}

	// Build member list
	members := make([]map[string]interface{}, 0, len(group.Members))
	for _, m := range group.Members {
		if m.User != nil {
			members = append(members, map[string]interface{}{
				"user_id":   m.UserID,
				"username":  m.User.Username,
				"nickname":  m.User.Nickname,
				"avatar":    m.User.Avatar,
				"role":      m.Role,
				"joined_at": m.JoinedAt,
			})
		}
	}

	return map[string]interface{}{
		"id":         group.ID,
		"name":       group.Name,
		"avatar":     group.Avatar,
		"owner_id":   group.OwnerID,
		"owner_name": group.Owner.Nickname,
		"created_at": group.CreatedAt,
		"members":    members,
	}, nil
}

// ============ group.join ============

type GroupJoinMethod struct {
	storage *storage.Storage
}

func NewGroupJoinMethod(s *storage.Storage) *GroupJoinMethod {
	return &GroupJoinMethod{storage: s}
}

func (m *GroupJoinMethod) Name() string { return "group.join" }

func (m *GroupJoinMethod) RequireAuth() bool { return true }

type GroupJoinParams struct {
	GroupID int64 `json:"group_id"`
}

func (m *GroupJoinMethod) Execute(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var p GroupJoinParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, fmt.Errorf("invalid params: %v", err)
	}

	if p.GroupID == 0 {
		return nil, errors.New("group_id is required")
	}

	userID := ctx.Value("user_id").(int64)
	db := m.storage.GetDB()

	// Check if group exists
	var group models.Group
	if err := db.First(&group, p.GroupID).Error; err != nil {
		return nil, errors.New("group not found")
	}

	// Check if already a member
	var existing models.GroupMember
	err := db.Where("group_id = ? AND user_id = ?", p.GroupID, userID).First(&existing).Error
	if err == nil {
		return nil, errors.New("already a member of this group")
	}

	member := &models.GroupMember{
		GroupID:  p.GroupID,
		UserID:   userID,
		Role:     models.GroupRoleMember,
		JoinedAt: time.Now(),
	}

	if err := db.Create(member).Error; err != nil {
		return nil, fmt.Errorf("failed to join group: %v", err)
	}

	return map[string]interface{}{
		"message": "joined group successfully",
	}, nil
}
