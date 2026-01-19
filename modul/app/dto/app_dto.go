package app_dto

type AppResponse struct{
	ID        int64  `json:"id"`
	AppPageID int64  `json:"app_page_id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Type      string `json:"type"`
	IsRequire bool   `json:"is_require"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}
