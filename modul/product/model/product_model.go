package product_model

import (
	"time"
)

type Product struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"type:varchar(200);not null;"`
	Slug      string    `json:"slug" gorm:"type:varchar(100);not null;"`
	Price     int64     `json:"price" gorm:"not null;"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at""`
}

func (Product) TableName() string {
	return "products"
}
