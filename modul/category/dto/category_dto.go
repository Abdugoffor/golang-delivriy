package category_dto

import (
	"my-project/helper"
	category_model "my-project/modul/category/model"
)

type Create struct {
	Name     string `json:"name" query:"name" form:"name"`
	IsActive bool   `json:"is_active" query:"is_active" form:"is_active"`
}

type Update struct {
	Name     string `json:"name" query:"name" form:"name"`
	IsActive bool   `json:"is_active" query:"is_active" form:"is_active"`
}

type Response struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type Filter struct {
	Name   string `json:"name" query:"name" form:"name"`
	Status string `json:"status" query:"status" form:"status"`
	Sort   string `json:"sort" query:"sort" form:"sort"`
	Column string `json:"column" query:"column" form:"column"`
}

func ToResponse(category category_model.Category) Response {
	return Response{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		IsActive:  category.IsActive,
		CreatedAt: helper.FormatDate(category.CreatedAt),
		UpdatedAt: helper.FormatDate(category.UpdatedAt),
		DeletedAt: helper.FormatDate(category.DeletedAt),
	}
}
