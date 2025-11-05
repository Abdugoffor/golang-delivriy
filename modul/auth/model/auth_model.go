package auth_model

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64          `json:"user_id" gorm:"index;not null"`
	Token     string         `json:"token" gorm:"type:text;not null;uniqueIndex"`
	ExpiresAt time.Time      `json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (Session) TableName() string {
	return "sessions"
}
