package app_model

import (
	"time"

	"gorm.io/gorm"
)

type AppCategory struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name" gorm:"type:varchar(200);not null;"`
	Slug      string         `json:"slug" gorm:"type:varchar(100);not null;"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (AppCategory) TableName() string {
	return "app_categories"
}
