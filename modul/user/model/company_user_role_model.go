package user_model

import (
	"time"

	"gorm.io/gorm"
)

type CompanyUserRole struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	CompanyID int64          `json:"company_id"`
	UserID    int64          `json:"user_id"`
	RoleID    int64          `json:"role_id"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	
	User User `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Role Role `json:"role" gorm:"foreignKey:RoleID;references:ID"`
}

func (CompanyUserRole) TableName() string {
	return "company_user_roles"
}
