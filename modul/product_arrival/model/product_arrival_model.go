package product_arrival_model

import (
	"time"
)

type ProductArrival struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement;"`
	Author    string    `json:"author" gorm:"not null;"`
	ProductID uint      `json:"product_id" gorm:"not null;"`
	Count     int       `json:"count" gorm:"not null;"`
	Sum       int       `json:"sum" gorm:"not null;"`
	CreatedAt time.Time `json:"created_at"`
}

func (ProductArrival) TableName() string {
	return "product_arrivals"
}
