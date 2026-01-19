package app_model

import (
	"time"

	"gorm.io/gorm"
)

type AppValue struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	AppID     int64          `json:"app_id"`
	AppPageID int64          `json:"app_page_id"`
	AppFormID int64          `json:"app_form_id"`
	Value     string         `json:"value"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
