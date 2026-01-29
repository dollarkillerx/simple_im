package ws

import (
	"testing"
	"time"
)

func TestNewHub(t *testing.T) {
	hub := NewHub()

	if hub == nil {
		t.Fatal("NewHub should not return nil")
	}

	if hub.clients == nil {
		t.Error("clients map should be initialized")
	}

	if hub.register == nil {
		t.Error("register channel should be initialized")
	}

	if hub.unregister == nil {
		t.Error("unregister channel should be initialized")
	}

	if hub.broadcast == nil {
		t.Error("broadcast channel should be initialized")
	}
}

func TestHub_IsOnline(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Initially no one is online
	if hub.IsOnline(1) {
		t.Error("User 1 should not be online initially")
	}

	// Manually add a client to test
	hub.mu.Lock()
	hub.clients[1] = &Client{UserID: 1}
	hub.mu.Unlock()

	if !hub.IsOnline(1) {
		t.Error("User 1 should be online after adding")
	}

	if hub.IsOnline(2) {
		t.Error("User 2 should not be online")
	}
}

func TestHub_RegisterUnregister(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{
		UserID:   100,
		Username: "testuser",
		send:     make(chan *Message, 256),
	}

	// Register
	hub.Register(client)
	time.Sleep(50 * time.Millisecond) // Wait for goroutine to process

	if !hub.IsOnline(100) {
		t.Error("User should be online after registration")
	}

	// Unregister
	hub.Unregister(client)
	time.Sleep(50 * time.Millisecond)

	if hub.IsOnline(100) {
		t.Error("User should be offline after unregistration")
	}
}

func TestHub_Broadcast_PrivateMessage(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create receiver client
	receiverChan := make(chan *Message, 10)
	receiver := &Client{
		UserID:   2,
		Username: "receiver",
		send:     receiverChan,
	}

	hub.Register(receiver)
	time.Sleep(50 * time.Millisecond)

	// Send private message
	msg := &Message{
		Type:       "message",
		SenderID:   1,
		ReceiverID: 2,
		Content:    "Hello!",
	}

	hub.Broadcast(msg)
	time.Sleep(50 * time.Millisecond)

	select {
	case received := <-receiverChan:
		if received.Content != "Hello!" {
			t.Errorf("Expected 'Hello!', got '%s'", received.Content)
		}
		if received.SenderID != 1 {
			t.Errorf("Expected sender 1, got %d", received.SenderID)
		}
	default:
		t.Error("Receiver should have received the message")
	}
}

func TestHub_Broadcast_GroupMessage(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create group members
	member1Chan := make(chan *Message, 10)
	member1 := &Client{UserID: 2, send: member1Chan}

	member2Chan := make(chan *Message, 10)
	member2 := &Client{UserID: 3, send: member2Chan}

	hub.Register(member1)
	hub.Register(member2)
	time.Sleep(50 * time.Millisecond)

	// Send group message (sender is user 1)
	msg := &Message{
		Type:         "message",
		SenderID:     1,
		GroupID:      100,
		Content:      "Group message",
		GroupMembers: []int64{1, 2, 3}, // All members including sender
	}

	hub.Broadcast(msg)
	time.Sleep(50 * time.Millisecond)

	// Member 1 should receive
	select {
	case received := <-member1Chan:
		if received.Content != "Group message" {
			t.Errorf("Member 1: Expected 'Group message', got '%s'", received.Content)
		}
	default:
		t.Error("Member 1 should have received the message")
	}

	// Member 2 should receive
	select {
	case received := <-member2Chan:
		if received.Content != "Group message" {
			t.Errorf("Member 2: Expected 'Group message', got '%s'", received.Content)
		}
	default:
		t.Error("Member 2 should have received the message")
	}
}

func TestMessage_Types(t *testing.T) {
	if MsgTypeText != 1 {
		t.Errorf("MsgTypeText should be 1, got %d", MsgTypeText)
	}
	if MsgTypeImage != 2 {
		t.Errorf("MsgTypeImage should be 2, got %d", MsgTypeImage)
	}
	if MsgTypeFile != 3 {
		t.Errorf("MsgTypeFile should be 3, got %d", MsgTypeFile)
	}
}
