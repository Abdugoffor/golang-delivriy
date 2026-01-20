package app_dto

import (
	"encoding/json"

	"github.com/Abdugoffor/echo-crud-pg/response"
)

type AppCatePage = response.PageData[AppCateResponse]

type AppCateResponse struct {
	ID        int64           `json:"id"`
	Name      string          `json:"name"`
	Slug      string          `json:"slug"`
	Pages     json.RawMessage `json:"pages"`
	IsActive  bool            `json:"is_active"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
	DeletedAt string          `json:"deleted_at"`
}
