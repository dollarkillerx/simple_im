package models

import "time"

type File struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    int64     `gorm:"not null;index" json:"user_id"`
	Filename  string    `gorm:"size:255;not null" json:"filename"`
	Filepath  string    `gorm:"size:500;not null" json:"filepath"`
	Filesize  int64     `gorm:"not null" json:"filesize"`
	Mimetype  string    `gorm:"size:100" json:"mimetype"`
	CreatedAt time.Time `json:"created_at"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (File) TableName() string {
	return "files"
}
