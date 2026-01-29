package jwt

import (
	"testing"
	"time"
)

func TestJWTManager_GenerateAndParseToken(t *testing.T) {
	manager := NewJWTManager("test_secret_key", 3600)

	userID := int64(123)
	username := "testuser"

	// Generate token
	token, err := manager.GenerateToken(userID, username)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Fatal("Token should not be empty")
	}

	// Parse token
	claims, err := manager.ParseToken(token)
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("UserID mismatch: expected %d, got %d", userID, claims.UserID)
	}

	if claims.Username != username {
		t.Errorf("Username mismatch: expected %s, got %s", username, claims.Username)
	}
}

func TestJWTManager_ExpiredToken(t *testing.T) {
	// Create manager with 1 second expiry
	manager := NewJWTManager("test_secret_key", 1)

	token, err := manager.GenerateToken(1, "testuser")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Wait for token to expire
	time.Sleep(2 * time.Second)

	_, err = manager.ParseToken(token)
	if err == nil {
		t.Fatal("Expected error for expired token")
	}

	if err != ErrTokenExpired {
		t.Errorf("Expected ErrTokenExpired, got: %v", err)
	}
}

func TestJWTManager_InvalidToken(t *testing.T) {
	manager := NewJWTManager("test_secret_key", 3600)

	// Test with invalid token
	_, err := manager.ParseToken("invalid.token.here")
	if err == nil {
		t.Fatal("Expected error for invalid token")
	}

	if err != ErrTokenInvalid {
		t.Errorf("Expected ErrTokenInvalid, got: %v", err)
	}
}

func TestJWTManager_WrongSecret(t *testing.T) {
	manager1 := NewJWTManager("secret1", 3600)
	manager2 := NewJWTManager("secret2", 3600)

	token, _ := manager1.GenerateToken(1, "user")

	_, err := manager2.ParseToken(token)
	if err == nil {
		t.Fatal("Expected error when parsing with wrong secret")
	}
}
