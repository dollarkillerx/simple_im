package models

import "time"

type GroupRole int

const (
	GroupRoleMember GroupRole = 0
	GroupRoleAdmin  GroupRole = 1
	GroupRoleOwner  GroupRole = 2
)

type Group struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	OwnerID   int64     `gorm:"not null" json:"owner_id"`
	Avatar    string    `gorm:"size:500" json:"avatar"`
	CreatedAt time.Time `json:"created_at"`

	Owner   *User          `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Members []GroupMember  `gorm:"foreignKey:GroupID" json:"members,omitempty"`
}

func (Group) TableName() string {
	return "groups"
}

type GroupMember struct {
	ID       int64     `gorm:"primaryKey" json:"id"`
	GroupID  int64     `gorm:"not null;uniqueIndex:idx_group_member" json:"group_id"`
	UserID   int64     `gorm:"not null;uniqueIndex:idx_group_member" json:"user_id"`
	Role     GroupRole `gorm:"default:0" json:"role"` // 0:member 1:admin 2:owner
	JoinedAt time.Time `json:"joined_at"`

	Group *Group `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	User  *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (GroupMember) TableName() string {
	return "group_members"
}
