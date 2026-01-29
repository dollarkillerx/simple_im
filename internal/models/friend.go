package models

import "time"

type FriendStatus int

const (
	FriendStatusPending  FriendStatus = 0
	FriendStatusAccepted FriendStatus = 1
	FriendStatusRejected FriendStatus = 2
)

type Friend struct {
	ID        int64        `gorm:"primaryKey" json:"id"`
	UserID    int64        `gorm:"not null;index:idx_friend_user" json:"user_id"`
	FriendID  int64        `gorm:"not null;index:idx_friend_user" json:"friend_id"`
	Status    FriendStatus `gorm:"default:0" json:"status"` // 0:pending 1:accepted 2:rejected
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`

	User   *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Friend *User `gorm:"foreignKey:FriendID" json:"friend,omitempty"`
}

func (Friend) TableName() string {
	return "friends"
}
