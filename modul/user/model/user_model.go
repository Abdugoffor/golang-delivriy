package user_model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                  uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name                string         `json:"name" gorm:"type:varchar(200);not null;"`
	Email               string         `json:"email" gorm:"type:varchar(200);uniqueIndex;not null;"`
	PasswordHash        string         `json:"-" gorm:"type:varchar(255);not null;"`
	IsVerified          bool           `json:"is_verified" gorm:"default:false"`
	VerificationToken   string         `json:"-" gorm:"type:varchar(255);index"`
	PasswordResetToken  string         `json:"-" gorm:"type:varchar(255);index"`
	PasswordResetExpiry *time.Time     `json:"-"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
