package history

import (
	"time"

	"gorm.io/gorm"
)

// History model
type History struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    *int64         `json:"user_id"`   // kim bajardi
	Table     *string        `json:"table"`     // qaysi jadval
	RowID     *string        `json:"row_id"`    // model PK
	Action    *string        `json:"action"`    // create/update/delete/restore
	OldValue  *string        `json:"old_value"` // eski qiymatlar (JSON)
	NewValue  *string        `json:"new_value"` // yangi qiymatlar (JSON)
	IP        *string        `json:"ip"`        // foydalanuvchi IP
	API       *string        `json:"api"`       // endpoint
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Table nomi
func (History) TableName() string {
	return "histories"
}
