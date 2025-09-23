package product_arrival_model

import (
	"time"
)

type ProductArrivalLog struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Type           string    `json:"type" gorm:"type:varchar(50);not null"` // arrival, sale, return, defect
	ProductID      uint      `json:"product_id" gorm:"not null;index"`
	QuantityBefore int       `json:"quantity_before" gorm:"not null"` // operatsiyadan oldingi qoldiq
	Quantity       int       `json:"quantity" gorm:"not null"`        // qancha keldi yoki chiqdi
	QuantityAfter  int       `json:"quantity_after" gorm:"not null"`  // operatsiyadan keyingi qoldiq
	Sum            int       `json:"sum" gorm:"not null"`             // summasi
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (ProductArrivalLog) TableName() string {
	return "product_arrival_logs"
}
