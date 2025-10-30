package product_dto

import (
	"my-project/helper"
	product_model "my-project/modul/product/model"
)

type Create struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type Update struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}
type Filter struct {
	Name  string `json:"name" query:"name" form:"name"`
	Price int    `json:"price" query:"price" form:"price"`
}

type Response struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Price     int    `json:"price"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

func ToResponse(product product_model.Product) Response {
	return Response{
		ID:        product.ID,
		Name:      product.Name,
		Slug:      product.Slug,
		Price:     product.Price,
		IsActive:  product.IsActive,
		CreatedAt: helper.FormatDate(product.CreatedAt),
		UpdatedAt: helper.FormatDate(product.UpdatedAt),
		DeletedAt: helper.FormatDate(product.DeletedAt.Time),
	}
}
