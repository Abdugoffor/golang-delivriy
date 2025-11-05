package user_model

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name" gorm:"type:varchar(200);not null;"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (Role) TableName() string {
	return "roles"
}
