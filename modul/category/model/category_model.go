package category_model

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name" gorm:"type:varchar(200);not null;"`
	Slug      string         `json:"slug" gorm:"type:varchar(100);not null;"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	Order     int            `json:"order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (Category) TableName() string {
	return "categories"
}
