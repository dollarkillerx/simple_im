package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"simple_im/internal/models"
	"simple_im/internal/storage"
	"simple_im/pkg/common/jwt"
)

// ============ user.register ============

type UserRegisterMethod struct {
	storage    *storage.Storage
	jwtManager *jwt.JWTManager
}

func NewUserRegisterMethod(s *storage.Storage, j *jwt.JWTManager) *UserRegisterMethod {
	return &UserRegisterMethod{storage: s, jwtManager: j}
}

func (m *UserRegisterMethod) Name() string { return "user.register" }

func (m *UserRegisterMethod) RequireAuth() bool { return false }

type UserRegisterParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

func (m *UserRegisterMethod) Execute(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var p UserRegisterParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, fmt.Errorf("invalid params: %v", err)
	}

	if p.Username == "" || p.Password == "" {
		return nil, errors.New("username and password are required")
	}

	if len(p.Username) < 3 || len(p.Username) > 50 {
		return nil, errors.New("username must be 3-50 characters")
	}

	if len(p.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	db := m.storage.GetDB()

	// Check if username exists
	var count int64
	db.Model(&models.User{}).Where("username = ?", p.Username).Count(&count)
	if count > 0 {
		return nil, errors.New("username already exists")
	}

	user := &models.User{
		Username: p.Username,
		Nickname: p.Nickname,
	}
	if user.Nickname == "" {
		user.Nickname = p.Username
	}

	if err := user.SetPassword(p.Password); err != nil {
		return nil, fmt.Errorf("failed to set password: %v", err)
	}

	if err := db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	token, err := m.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return map[string]interface{}{
		"user":  user,
		"token": token,
	}, nil
}

// ============ user.login ============

type UserLoginMethod struct {
	storage    *storage.Storage
	jwtManager *jwt.JWTManager
}

func NewUserLoginMethod(s *storage.Storage, j *jwt.JWTManager) *UserLoginMethod {
	return &UserLoginMethod{storage: s, jwtManager: j}
}

func (m *UserLoginMethod) Name() string { return "user.login" }

func (m *UserLoginMethod) RequireAuth() bool { return false }

type UserLoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (m *UserLoginMethod) Execute(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var p UserLoginParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, fmt.Errorf("invalid params: %v", err)
	}

	if p.Username == "" || p.Password == "" {
		return nil, errors.New("username and password are required")
	}

	db := m.storage.GetDB()

	var user models.User
	if err := db.Where("username = ?", p.Username).First(&user).Error; err != nil {
		return nil, errors.New("invalid username or password")
	}

	if !user.CheckPassword(p.Password) {
		return nil, errors.New("invalid username or password")
	}

	if user.Status != 1 {
		return nil, errors.New("user is disabled")
	}

	token, err := m.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return map[string]interface{}{
		"user":  user,
		"token": token,
	}, nil
}

// ============ user.info ============

type UserInfoMethod struct {
	storage *storage.Storage
}

func NewUserInfoMethod(s *storage.Storage) *UserInfoMethod {
	return &UserInfoMethod{storage: s}
}

func (m *UserInfoMethod) Name() string { return "user.info" }

func (m *UserInfoMethod) RequireAuth() bool { return true }

type UserInfoParams struct {
	UserID int64 `json:"user_id"`
}

func (m *UserInfoMethod) Execute(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var p UserInfoParams
	if len(params) > 0 {
		json.Unmarshal(params, &p)
	}

	// If no user_id specified, return current user info
	userID := p.UserID
	if userID == 0 {
		if v := ctx.Value("user_id"); v != nil {
			userID = v.(int64)
		}
	}

	if userID == 0 {
		return nil, errors.New("user_id is required")
	}

	db := m.storage.GetDB()

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
