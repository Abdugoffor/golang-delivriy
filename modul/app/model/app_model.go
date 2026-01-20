package app_model

import (
	"time"

	"gorm.io/gorm"
)

type App struct {
	ID            int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID        int64          `json:"user_id"`
	AppCategoryID int64          `json:"app_category_id"`
	IsActive      bool           `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (App) TableName() string {
	return "apps"
}
