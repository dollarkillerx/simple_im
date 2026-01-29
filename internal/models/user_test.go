package models

import (
	"testing"
)

func TestUser_SetPassword(t *testing.T) {
	user := &User{}
	password := "testpassword123"

	err := user.SetPassword(password)
	if err != nil {
		t.Fatalf("Failed to set password: %v", err)
	}

	if user.Password == "" {
		t.Fatal("Password should not be empty after setting")
	}

	if user.Password == password {
		t.Fatal("Password should be hashed, not plain text")
	}
}

func TestUser_CheckPassword(t *testing.T) {
	user := &User{}
	password := "testpassword123"

	err := user.SetPassword(password)
	if err != nil {
		t.Fatalf("Failed to set password: %v", err)
	}

	// Test correct password
	if !user.CheckPassword(password) {
		t.Error("CheckPassword should return true for correct password")
	}

	// Test wrong password
	if user.CheckPassword("wrongpassword") {
		t.Error("CheckPassword should return false for wrong password")
	}
}

func TestUser_TableName(t *testing.T) {
	user := User{}
	if user.TableName() != "users" {
		t.Errorf("Expected table name 'users', got '%s'", user.TableName())
	}
}

func TestFriend_TableName(t *testing.T) {
	friend := Friend{}
	if friend.TableName() != "friends" {
		t.Errorf("Expected table name 'friends', got '%s'", friend.TableName())
	}
}

func TestGroup_TableName(t *testing.T) {
	group := Group{}
	if group.TableName() != "groups" {
		t.Errorf("Expected table name 'groups', got '%s'", group.TableName())
	}
}

func TestGroupMember_TableName(t *testing.T) {
	member := GroupMember{}
	if member.TableName() != "group_members" {
		t.Errorf("Expected table name 'group_members', got '%s'", member.TableName())
	}
}

func TestMessage_TableName(t *testing.T) {
	msg := Message{}
	if msg.TableName() != "messages" {
		t.Errorf("Expected table name 'messages', got '%s'", msg.TableName())
	}
}

func TestFile_TableName(t *testing.T) {
	file := File{}
	if file.TableName() != "files" {
		t.Errorf("Expected table name 'files', got '%s'", file.TableName())
	}
}
