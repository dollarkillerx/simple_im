package ws

import "time"

type MessageType int

const (
	MsgTypeText  MessageType = 1
	MsgTypeImage MessageType = 2
	MsgTypeFile  MessageType = 3
)

type Message struct {
	ID           int64       `json:"id"`
	Type         string      `json:"type"` // "message", "notification", "friend_request", etc.
	SenderID     int64       `json:"sender_id"`
	SenderName   string      `json:"sender_name,omitempty"`
	ReceiverID   int64       `json:"receiver_id,omitempty"`
	GroupID      int64       `json:"group_id,omitempty"`
	GroupName    string      `json:"group_name,omitempty"`
	MsgType      MessageType `json:"msg_type"`
	Content      string      `json:"content,omitempty"`
	FileURL      string      `json:"file_url,omitempty"`
	FileName     string      `json:"file_name,omitempty"`
	FileSize     int64       `json:"file_size,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	GroupMembers []int64     `json:"-"` // Internal use for broadcasting
}
