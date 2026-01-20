package app_dto

import (
	"encoding/json"

	"github.com/Abdugoffor/echo-crud-pg/response"
)

type AppPage = response.PageData[AppResponse]

type AppResponse struct {
	ID          int64           `json:"id"`
	UserID      int64           `json:"user_id"`
	AppCategory json.RawMessage `json:"app_category"`
	AppPage     json.RawMessage `json:"app_page"`
	IsActive    bool            `json:"is_active"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
	DeletedAt   string          `json:"deleted_at"`
}
