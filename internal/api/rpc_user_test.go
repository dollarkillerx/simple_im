package api

import (
	"context"
	"encoding/json"
	"simple_im/internal/models"
	"testing"
)

func TestUserRegisterMethod_Execute(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	method := NewUserRegisterMethod(env.Storage, env.JWTManager)

	// Test successful registration
	params, _ := json.Marshal(UserRegisterParams{
		Username: "newuser",
		Password: "password123",
		Nickname: "New User",
	})

	result, err := method.Execute(context.Background(), params)
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}

	resultMap := result.(map[string]interface{})
	if resultMap["token"] == "" {
		t.Error("Token should not be empty")
	}

	// Test duplicate username
	_, err = method.Execute(context.Background(), params)
	if err == nil {
		t.Error("Should fail for duplicate username")
	}
}

func TestUserRegisterMethod_Validation(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	method := NewUserRegisterMethod(env.Storage, env.JWTManager)

	tests := []struct {
		name    string
		params  UserRegisterParams
		wantErr bool
	}{
		{
			name:    "Empty username",
			params:  UserRegisterParams{Username: "", Password: "password123"},
			wantErr: true,
		},
		{
			name:    "Empty password",
			params:  UserRegisterParams{Username: "testuser", Password: ""},
			wantErr: true,
		},
		{
			name:    "Short username",
			params:  UserRegisterParams{Username: "ab", Password: "password123"},
			wantErr: true,
		},
		{
			name:    "Short password",
			params:  UserRegisterParams{Username: "testuser", Password: "12345"},
			wantErr: true,
		},
		{
			name:    "Valid registration",
			params:  UserRegisterParams{Username: "validuser", Password: "password123"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, _ := json.Marshal(tt.params)
			_, err := method.Execute(context.Background(), params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserLoginMethod_Execute(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	// Create a test user first
	_, err = env.CreateTestUser("loginuser", "password123")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	method := NewUserLoginMethod(env.Storage, env.JWTManager)

	// Test successful login
	params, _ := json.Marshal(UserLoginParams{
		Username: "loginuser",
		Password: "password123",
	})

	result, err := method.Execute(context.Background(), params)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	resultMap := result.(map[string]interface{})
	if resultMap["token"] == "" {
		t.Error("Token should not be empty")
	}

	// Test wrong password
	params, _ = json.Marshal(UserLoginParams{
		Username: "loginuser",
		Password: "wrongpassword",
	})

	_, err = method.Execute(context.Background(), params)
	if err == nil {
		t.Error("Should fail for wrong password")
	}

	// Test non-existent user
	params, _ = json.Marshal(UserLoginParams{
		Username: "nonexistent",
		Password: "password123",
	})

	_, err = method.Execute(context.Background(), params)
	if err == nil {
		t.Error("Should fail for non-existent user")
	}
}

func TestUserInfoMethod_Execute(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user, err := env.CreateTestUser("infouser", "password123")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	method := NewUserInfoMethod(env.Storage)

	// Test with user_id in context
	ctx := context.WithValue(context.Background(), "user_id", user.ID)

	result, err := method.Execute(ctx, nil)
	if err != nil {
		t.Fatalf("GetInfo failed: %v", err)
	}

	userResult := result.(models.User)
	if userResult.Username != "infouser" {
		t.Errorf("Expected username 'infouser', got '%s'", userResult.Username)
	}

	// Test with explicit user_id in params
	params, _ := json.Marshal(UserInfoParams{UserID: user.ID})
	result, err = method.Execute(context.Background(), params)
	if err != nil {
		t.Fatalf("GetInfo with params failed: %v", err)
	}
}

func TestUserMethods_RequireAuth(t *testing.T) {
	env, _ := SetupTestEnv()

	registerMethod := NewUserRegisterMethod(env.Storage, env.JWTManager)
	loginMethod := NewUserLoginMethod(env.Storage, env.JWTManager)
	infoMethod := NewUserInfoMethod(env.Storage)

	if registerMethod.RequireAuth() {
		t.Error("Register should not require auth")
	}

	if loginMethod.RequireAuth() {
		t.Error("Login should not require auth")
	}

	if !infoMethod.RequireAuth() {
		t.Error("Info should require auth")
	}
}

func TestUserMethods_Name(t *testing.T) {
	env, _ := SetupTestEnv()

	registerMethod := NewUserRegisterMethod(env.Storage, env.JWTManager)
	loginMethod := NewUserLoginMethod(env.Storage, env.JWTManager)
	infoMethod := NewUserInfoMethod(env.Storage)

	if registerMethod.Name() != "user.register" {
		t.Errorf("Expected 'user.register', got '%s'", registerMethod.Name())
	}

	if loginMethod.Name() != "user.login" {
		t.Errorf("Expected 'user.login', got '%s'", loginMethod.Name())
	}

	if infoMethod.Name() != "user.info" {
		t.Errorf("Expected 'user.info', got '%s'", infoMethod.Name())
	}
}
