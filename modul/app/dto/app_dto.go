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

type Create struct {
	AppCategoryID int64 `json:"app_category_id" validate:"required"`
}
type CreateCate struct {
	AppID        int64           `json:"app_id" validate:"required"` // App / Category ID
	AppFormField []AppCreateForm `json:"forms" validate:"required"`  // formalar ro'yxati
}

type AppCreateForm struct {
	AppPageID int64  `json:"app_page_id" validate:"required"` // qaysi page
	AppFormID int64  `json:"app_form_id" validate:"required"` // qaysi field (forma)
	Value     string `json:"value" validate:"required"`       // foydalanuvchi kiritgan qiymat
}
