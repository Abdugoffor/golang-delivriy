package company_model

import (
	category_model "my-project/modul/category/model"
	user_model "my-project/modul/user/model"
	"time"

	"gorm.io/gorm"
)

type Company struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	ParentID  int64          `json:"parent_id"`
	UserID    int64          `json:"user_id"`
	Name      string         `json:"name" gorm:"type:varchar(200);not null;"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	Categories       []category_model.Category    `json:"categories" gorm:"foreignKey:CompanyID;references:ID"`                                                                 // ✅
	Users            []user_model.User            `json:"users" gorm:"many2many:company_user_roles;foreignKey:ID;joinForeignKey:CompanyID;References:ID;joinReferences:UserID"` // ✅
	CompanyUserRoles []user_model.CompanyUserRole `json:"company_user_roles" gorm:"foreignKey:CompanyID;references:ID"`
	
}

func (Company) TableName() string {
	return "company"
}
