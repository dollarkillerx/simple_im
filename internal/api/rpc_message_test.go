package api

import (
	"context"
	"encoding/json"
	"simple_im/internal/models"
	"testing"
	"time"
)

func TestMessageSendMethod_PrivateMessage(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user1, _ := env.CreateTestUser("sender", "password")
	user2, _ := env.CreateTestUser("receiver", "password")

	// Create friendship
	env.CreateTestFriendship(user1.ID, user2.ID, models.FriendStatusAccepted)

	method := NewMessageSendMethod(env.Storage, env.Hub)

	ctx := context.WithValue(context.Background(), "user_id", user1.ID)
	ctx = context.WithValue(ctx, "username", user1.Username)

	params, _ := json.Marshal(MessageSendParams{
		ReceiverID: user2.ID,
		MsgType:    models.MsgTypeText,
		Content:    "Hello!",
	})

	result, err := method.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Send message failed: %v", err)
	}

	msg := result.(*models.Message)
	if msg.Content != "Hello!" {
		t.Errorf("Expected content 'Hello!', got '%s'", msg.Content)
	}
	if msg.SenderID != user1.ID {
		t.Errorf("Expected sender %d, got %d", user1.ID, msg.SenderID)
	}
	if msg.ReceiverID == nil || *msg.ReceiverID != user2.ID {
		t.Errorf("Expected receiver %d, got %v", user2.ID, msg.ReceiverID)
	}
}

func TestMessageSendMethod_GroupMessage(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user1, _ := env.CreateTestUser("groupsender", "password")
	group, _ := env.CreateTestGroup("Message Group", user1.ID)

	method := NewMessageSendMethod(env.Storage, env.Hub)

	ctx := context.WithValue(context.Background(), "user_id", user1.ID)
	ctx = context.WithValue(ctx, "username", user1.Username)

	params, _ := json.Marshal(MessageSendParams{
		GroupID: group.ID,
		MsgType: models.MsgTypeText,
		Content: "Hello group!",
	})

	result, err := method.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Send group message failed: %v", err)
	}

	msg := result.(*models.Message)
	if msg.GroupID == nil || *msg.GroupID != group.ID {
		t.Errorf("Expected group ID %d, got %v", group.ID, msg.GroupID)
	}
}

func TestMessageSendMethod_NotFriend(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user1, _ := env.CreateTestUser("notfriend1", "password")
	user2, _ := env.CreateTestUser("notfriend2", "password")

	// No friendship created

	method := NewMessageSendMethod(env.Storage, env.Hub)

	ctx := context.WithValue(context.Background(), "user_id", user1.ID)
	ctx = context.WithValue(ctx, "username", user1.Username)

	params, _ := json.Marshal(MessageSendParams{
		ReceiverID: user2.ID,
		MsgType:    models.MsgTypeText,
		Content:    "Hello!",
	})

	_, err = method.Execute(ctx, params)
	if err == nil {
		t.Error("Should not be able to send message to non-friend")
	}
}

func TestMessageSendMethod_NotGroupMember(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user1, _ := env.CreateTestUser("groupowner4", "password")
	user2, _ := env.CreateTestUser("notmember", "password")
	group, _ := env.CreateTestGroup("Exclusive Group", user1.ID)

	method := NewMessageSendMethod(env.Storage, env.Hub)

	ctx := context.WithValue(context.Background(), "user_id", user2.ID)
	ctx = context.WithValue(ctx, "username", user2.Username)

	params, _ := json.Marshal(MessageSendParams{
		GroupID: group.ID,
		MsgType: models.MsgTypeText,
		Content: "Hello!",
	})

	_, err = method.Execute(ctx, params)
	if err == nil {
		t.Error("Non-member should not be able to send message to group")
	}
}

func TestMessageSendMethod_ImageMessage(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user1, _ := env.CreateTestUser("imgsender", "password")
	user2, _ := env.CreateTestUser("imgreceiver", "password")
	env.CreateTestFriendship(user1.ID, user2.ID, models.FriendStatusAccepted)

	method := NewMessageSendMethod(env.Storage, env.Hub)

	ctx := context.WithValue(context.Background(), "user_id", user1.ID)
	ctx = context.WithValue(ctx, "username", user1.Username)

	params, _ := json.Marshal(MessageSendParams{
		ReceiverID: user2.ID,
		MsgType:    models.MsgTypeImage,
		FileURL:    "/files/test.jpg",
		FileName:   "test.jpg",
		FileSize:   1024,
	})

	result, err := method.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Send image message failed: %v", err)
	}

	msg := result.(*models.Message)
	if msg.MsgType != models.MsgTypeImage {
		t.Errorf("Expected image type, got %d", msg.MsgType)
	}
	if msg.FileURL != "/files/test.jpg" {
		t.Errorf("Expected file URL '/files/test.jpg', got '%s'", msg.FileURL)
	}
}

func TestMessageSendMethod_Validation(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user, _ := env.CreateTestUser("validator", "password")
	method := NewMessageSendMethod(env.Storage, env.Hub)

	ctx := context.WithValue(context.Background(), "user_id", user.ID)
	ctx = context.WithValue(ctx, "username", user.Username)

	tests := []struct {
		name    string
		params  MessageSendParams
		wantErr bool
	}{
		{
			name:    "No receiver or group",
			params:  MessageSendParams{Content: "Hello"},
			wantErr: true,
		},
		{
			name:    "Text message without content",
			params:  MessageSendParams{ReceiverID: 1, MsgType: models.MsgTypeText},
			wantErr: true,
		},
		{
			name:    "Image message without file_url",
			params:  MessageSendParams{ReceiverID: 1, MsgType: models.MsgTypeImage},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, _ := json.Marshal(tt.params)
			_, err := method.Execute(ctx, params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessageHistoryMethod_PrivateChat(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user1, _ := env.CreateTestUser("hist1", "password")
	user2, _ := env.CreateTestUser("hist2", "password")
	env.CreateTestFriendship(user1.ID, user2.ID, models.FriendStatusAccepted)

	// Create some messages
	for i := 0; i < 5; i++ {
		receiverID := user2.ID
		msg := &models.Message{
			SenderID:   user1.ID,
			ReceiverID: &receiverID,
			MsgType:    models.MsgTypeText,
			Content:    "Message " + string(rune('A'+i)),
			CreatedAt:  time.Now().Add(time.Duration(i) * time.Second),
		}
		env.DB.Create(msg)
	}

	method := NewMessageHistoryMethod(env.Storage)

	ctx := context.WithValue(context.Background(), "user_id", user1.ID)

	params, _ := json.Marshal(MessageHistoryParams{
		ReceiverID: user2.ID,
		Limit:      10,
	})

	result, err := method.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Get message history failed: %v", err)
	}

	messages := result.([]models.Message)
	if len(messages) != 5 {
		t.Errorf("Expected 5 messages, got %d", len(messages))
	}
}

func TestMessageHistoryMethod_GroupChat(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user, _ := env.CreateTestUser("grouphist", "password")
	group, _ := env.CreateTestGroup("History Group", user.ID)

	// Create some group messages
	for i := 0; i < 3; i++ {
		groupID := group.ID
		msg := &models.Message{
			SenderID:  user.ID,
			GroupID:   &groupID,
			MsgType:   models.MsgTypeText,
			Content:   "Group message " + string(rune('A'+i)),
			CreatedAt: time.Now().Add(time.Duration(i) * time.Second),
		}
		env.DB.Create(msg)
	}

	method := NewMessageHistoryMethod(env.Storage)

	ctx := context.WithValue(context.Background(), "user_id", user.ID)

	params, _ := json.Marshal(MessageHistoryParams{
		GroupID: group.ID,
		Limit:   10,
	})

	result, err := method.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Get group message history failed: %v", err)
	}

	messages := result.([]models.Message)
	if len(messages) != 3 {
		t.Errorf("Expected 3 messages, got %d", len(messages))
	}
}

func TestMessageHistoryMethod_Pagination(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user1, _ := env.CreateTestUser("page1", "password")
	user2, _ := env.CreateTestUser("page2", "password")
	env.CreateTestFriendship(user1.ID, user2.ID, models.FriendStatusAccepted)

	// Create 10 messages
	var lastID int64
	for i := 0; i < 10; i++ {
		receiverID := user2.ID
		msg := &models.Message{
			SenderID:   user1.ID,
			ReceiverID: &receiverID,
			MsgType:    models.MsgTypeText,
			Content:    "Message " + string(rune('A'+i)),
			CreatedAt:  time.Now().Add(time.Duration(i) * time.Second),
		}
		env.DB.Create(msg)
		lastID = msg.ID
	}

	method := NewMessageHistoryMethod(env.Storage)
	ctx := context.WithValue(context.Background(), "user_id", user1.ID)

	// Get first 5 messages
	params, _ := json.Marshal(MessageHistoryParams{
		ReceiverID: user2.ID,
		Limit:      5,
	})

	result, _ := method.Execute(ctx, params)
	messages := result.([]models.Message)
	if len(messages) != 5 {
		t.Errorf("Expected 5 messages, got %d", len(messages))
	}

	// Get older messages using before_id
	params, _ = json.Marshal(MessageHistoryParams{
		ReceiverID: user2.ID,
		BeforeID:   lastID - 4, // Get messages before the 6th message
		Limit:      5,
	})

	result, _ = method.Execute(ctx, params)
	messages = result.([]models.Message)
	if len(messages) != 5 {
		t.Errorf("Expected 5 older messages, got %d", len(messages))
	}
}

func TestMessageHistoryMethod_NotGroupMember(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user1, _ := env.CreateTestUser("groupowner5", "password")
	user2, _ := env.CreateTestUser("outsider2", "password")
	group, _ := env.CreateTestGroup("Private History", user1.ID)

	method := NewMessageHistoryMethod(env.Storage)

	ctx := context.WithValue(context.Background(), "user_id", user2.ID)

	params, _ := json.Marshal(MessageHistoryParams{
		GroupID: group.ID,
	})

	_, err = method.Execute(ctx, params)
	if err == nil {
		t.Error("Non-member should not be able to view group history")
	}
}

func TestMessageMethods_RequireAuth(t *testing.T) {
	env, _ := SetupTestEnv()

	sendMethod := NewMessageSendMethod(env.Storage, env.Hub)
	historyMethod := NewMessageHistoryMethod(env.Storage)

	if !sendMethod.RequireAuth() {
		t.Error("Send should require auth")
	}
	if !historyMethod.RequireAuth() {
		t.Error("History should require auth")
	}
}

func TestMessageMethods_Name(t *testing.T) {
	env, _ := SetupTestEnv()

	sendMethod := NewMessageSendMethod(env.Storage, env.Hub)
	historyMethod := NewMessageHistoryMethod(env.Storage)

	if sendMethod.Name() != "message.send" {
		t.Errorf("Expected 'message.send', got '%s'", sendMethod.Name())
	}
	if historyMethod.Name() != "message.history" {
		t.Errorf("Expected 'message.history', got '%s'", historyMethod.Name())
	}
}
