package models

import (
	"time"

	"gorm.io/gorm"
)

type ChatMessage struct {
	ID        string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    *string        `gorm:"type:uuid;index" json:"user_id,omitempty"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	SessionID string         `gorm:"type:uuid;index" json:"session_id"`
	Message   string         `gorm:"type:text;not null" json:"message"`
	IsUser    bool           `gorm:"not null" json:"is_user"` // true for user message, false for AI response
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ChatMessage) TableName() string {
	return "chat_messages"
}
