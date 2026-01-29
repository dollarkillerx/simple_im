package api

import (
	"context"
	"encoding/json"
	"simple_im/internal/models"
	"testing"
)

func TestFriendAddMethod_Execute(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user1, _ := env.CreateTestUser("user1", "password")
	user2, _ := env.CreateTestUser("user2", "password")

	method := NewFriendAddMethod(env.Storage, env.Hub)

	ctx := context.WithValue(context.Background(), "user_id", user1.ID)
	ctx = context.WithValue(ctx, "username", user1.Username)

	// Test add friend by ID
	params, _ := json.Marshal(FriendAddParams{FriendID: user2.ID})
	result, err := method.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Add friend failed: %v", err)
	}

	resultMap := result.(map[string]interface{})
	if resultMap["message"] != "friend request sent" {
		t.Error("Expected success message")
	}

	// Verify friend request was created
	var friend models.Friend
	env.DB.Where("user_id = ? AND friend_id = ?", user1.ID, user2.ID).First(&friend)
	if friend.Status != models.FriendStatusPending {
		t.Error("Friend request should be pending")
	}
}

func TestFriendAddMethod_AddByUsername(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user1, _ := env.CreateTestUser("adder", "password")
	user2, _ := env.CreateTestUser("addee", "password")

	method := NewFriendAddMethod(env.Storage, env.Hub)

	ctx := context.WithValue(context.Background(), "user_id", user1.ID)
	ctx = context.WithValue(ctx, "username", user1.Username)

	params, _ := json.Marshal(FriendAddParams{Username: user2.Username})
	_, err = method.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Add friend by username failed: %v", err)
	}
}

func TestFriendAddMethod_CannotAddSelf(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user, _ := env.CreateTestUser("selfadder", "password")
	method := NewFriendAddMethod(env.Storage, env.Hub)

	ctx := context.WithValue(context.Background(), "user_id", user.ID)
	ctx = context.WithValue(ctx, "username", user.Username)

	params, _ := json.Marshal(FriendAddParams{FriendID: user.ID})
	_, err = method.Execute(ctx, params)
	if err == nil {
		t.Error("Should not be able to add self as friend")
	}
}

func TestFriendAcceptMethod_Execute(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user1, _ := env.CreateTestUser("requester", "password")
	user2, _ := env.CreateTestUser("accepter", "password")

	// Create pending friend request
	friend := &models.Friend{
		UserID:   user1.ID,
		FriendID: user2.ID,
		Status:   models.FriendStatusPending,
	}
	env.DB.Create(friend)

	method := NewFriendAcceptMethod(env.Storage, env.Hub)

	ctx := context.WithValue(context.Background(), "user_id", user2.ID)
	ctx = context.WithValue(ctx, "username", user2.Username)

	// Accept the request
	params, _ := json.Marshal(FriendAcceptParams{RequestID: friend.ID, Accept: true})
	_, err = method.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Accept friend failed: %v", err)
	}

	// Verify status changed
	var updated models.Friend
	env.DB.First(&updated, friend.ID)
	if updated.Status != models.FriendStatusAccepted {
		t.Error("Friend status should be accepted")
	}
}

func TestFriendAcceptMethod_Reject(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user1, _ := env.CreateTestUser("req2", "password")
	user2, _ := env.CreateTestUser("rej2", "password")

	friend := &models.Friend{
		UserID:   user1.ID,
		FriendID: user2.ID,
		Status:   models.FriendStatusPending,
	}
	env.DB.Create(friend)

	method := NewFriendAcceptMethod(env.Storage, env.Hub)

	ctx := context.WithValue(context.Background(), "user_id", user2.ID)
	ctx = context.WithValue(ctx, "username", user2.Username)

	params, _ := json.Marshal(FriendAcceptParams{RequestID: friend.ID, Accept: false})
	_, err = method.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Reject friend failed: %v", err)
	}

	var updated models.Friend
	env.DB.First(&updated, friend.ID)
	if updated.Status != models.FriendStatusRejected {
		t.Error("Friend status should be rejected")
	}
}

func TestFriendListMethod_Execute(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user1, _ := env.CreateTestUser("lister", "password")
	user2, _ := env.CreateTestUser("friend1", "password")
	user3, _ := env.CreateTestUser("friend2", "password")

	// Create accepted friendships
	env.CreateTestFriendship(user1.ID, user2.ID, models.FriendStatusAccepted)
	env.CreateTestFriendship(user3.ID, user1.ID, models.FriendStatusAccepted)

	method := NewFriendListMethod(env.Storage)

	ctx := context.WithValue(context.Background(), "user_id", user1.ID)

	result, err := method.Execute(ctx, nil)
	if err != nil {
		t.Fatalf("List friends failed: %v", err)
	}

	friends := result.([]map[string]interface{})
	if len(friends) != 2 {
		t.Errorf("Expected 2 friends, got %d", len(friends))
	}
}

func TestFriendPendingMethod_Execute(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user1, _ := env.CreateTestUser("pending1", "password")
	user2, _ := env.CreateTestUser("pending2", "password")

	// Create pending request where user2 is the receiver
	friend := &models.Friend{
		UserID:   user1.ID,
		FriendID: user2.ID,
		Status:   models.FriendStatusPending,
	}
	env.DB.Create(friend)

	method := NewFriendPendingMethod(env.Storage)

	ctx := context.WithValue(context.Background(), "user_id", user2.ID)

	result, err := method.Execute(ctx, nil)
	if err != nil {
		t.Fatalf("Get pending failed: %v", err)
	}

	pending := result.([]map[string]interface{})
	if len(pending) != 1 {
		t.Errorf("Expected 1 pending request, got %d", len(pending))
	}
}

func TestFriendMethods_RequireAuth(t *testing.T) {
	env, _ := SetupTestEnv()

	listMethod := NewFriendListMethod(env.Storage)
	addMethod := NewFriendAddMethod(env.Storage, env.Hub)
	acceptMethod := NewFriendAcceptMethod(env.Storage, env.Hub)
	pendingMethod := NewFriendPendingMethod(env.Storage)

	if !listMethod.RequireAuth() {
		t.Error("List should require auth")
	}
	if !addMethod.RequireAuth() {
		t.Error("Add should require auth")
	}
	if !acceptMethod.RequireAuth() {
		t.Error("Accept should require auth")
	}
	if !pendingMethod.RequireAuth() {
		t.Error("Pending should require auth")
	}
}
