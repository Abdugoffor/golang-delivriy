package history_model

import "time"

type History struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64     `json:"user_id" gorm:"not null;"`
	TableName string    `json:"table_name" gorm:"not null;"`
	RowID     int64     `json:"row_id" gorm:"not null;"`
	Action    string    `json:"action" gorm:"not null;"`
	OldValue  string    `json:"old_value" gorm:"not null;"`
	NewValue  string    `json:"new_value" gorm:"not null;"`
	IP        string    `json:"ip" gorm:"not null;"`
	API       string    `json:"api" gorm:"not null;"`
	Method    string    `json:"method" gorm:"not null;"`
	CreatedAt time.Time `json:"created_at"`
}
