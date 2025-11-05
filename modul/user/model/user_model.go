package user_model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name" gorm:"type:varchar(200);not null;"`
	Email     string         `json:"email" gorm:"type:varchar(200);not null;"`
	Password  string         `json:"password" gorm:"type:varchar(200);not null;"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
