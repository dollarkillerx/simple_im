package models

import "time"

type MessageType int

const (
	MsgTypeText  MessageType = 1
	MsgTypeImage MessageType = 2
	MsgTypeFile  MessageType = 3
)

type Message struct {
	ID         int64       `gorm:"primaryKey" json:"id"`
	SenderID   int64       `gorm:"not null;index" json:"sender_id"`
	ReceiverID *int64      `gorm:"index" json:"receiver_id,omitempty"` // Private chat (nullable)
	GroupID    *int64      `gorm:"index" json:"group_id,omitempty"`    // Group chat (nullable)
	MsgType    MessageType `gorm:"not null" json:"msg_type"`           // 1:text 2:image 3:file
	Content    string      `gorm:"type:text" json:"content,omitempty"`
	FileURL    string      `gorm:"size:500" json:"file_url,omitempty"`
	FileName   string      `gorm:"size:255" json:"file_name,omitempty"`
	FileSize   int64       `json:"file_size,omitempty"`
	CreatedAt  time.Time   `gorm:"index" json:"created_at"`

	Sender   *User  `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
	Receiver *User  `gorm:"foreignKey:ReceiverID;constraint:OnDelete:SET NULL" json:"receiver,omitempty"`
	Group    *Group `gorm:"foreignKey:GroupID;constraint:OnDelete:SET NULL" json:"group,omitempty"`
}

func (Message) TableName() string {
	return "messages"
}
