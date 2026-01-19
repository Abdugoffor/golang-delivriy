package app_model

import (
	"time"

	"gorm.io/gorm"
)

type AppForm struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	AppPageID int64          `json:"app_page_id"`
	Name      string         `json:"name" gorm:"type:varchar(200);not null;"`
	Slug      string         `json:"slug" gorm:"type:varchar(100);not null;"`
	Type      string         `json:"type" gorm:"type:varchar(100);not null;"`
	IsRequire bool           `json:"is_require" gorm:"default:false"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (AppForm) TableName() string {
	return "app_form"
}
