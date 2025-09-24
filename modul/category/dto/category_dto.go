package category_dto

import (
	"my-project/helper"
	category_model "my-project/modul/category/model"
)

type Create struct {
	Name     string `json:"name"`
	Order    int    `json:"order"`
	IsActive bool   `json:"is_active"`
}

type Update struct {
	Name     string `json:"name"`
	Order    int    `json:"order"`
	IsActive bool   `json:"is_active"`
}

type Response struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	IsActive  bool   `json:"is_active"`
	Order     int    `json:"order"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

func ToResponse(category category_model.Category) Response {
	return Response{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		IsActive:  category.IsActive,
		Order:     category.Order,
		CreatedAt: helper.FormatDate(category.CreatedAt),
		UpdatedAt: helper.FormatDate(category.UpdatedAt),
		DeletedAt: helper.FormatDate(category.DeletedAt.Time),
	}
}
