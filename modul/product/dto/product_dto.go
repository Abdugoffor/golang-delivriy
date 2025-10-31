package product_dto

import (
	"my-project/helper"
	product_model "my-project/modul/product/model"
)

type Create struct {
	Name  string `json:"name" query:"name" form:"name"`
	Price int64  `json:"price" query:"price" form:"price"`
}

type Update struct {
	Name  string `json:"name" query:"name" form:"name"`
	Price int64  `json:"price" query:"price" form:"price"`
}
type Filter struct {
	Name   string `json:"name" query:"name" form:"name"`
	Price  int64  `json:"price" query:"price" form:"price"`
	Status string `json:"status" query:"status" form:"status"`
	Sort   string `json:"sort" query:"sort" form:"sort"`
	Column string `json:"column" query:"column" form:"column"`
}

type Response struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Price     int64  `json:"price"`
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
		DeletedAt: helper.FormatDate(product.DeletedAt),
	}
}
