package company_dto

import (
	"my-project/helper"
	company_model "my-project/modul/company/model"
)

type Filter struct {
	Name   string `json:"name" query:"name" form:"name"`
	Status string `json:"status" query:"status" form:"status"`
	Sort   string `json:"sort" query:"sort" form:"sort"`
	Column string `json:"column" query:"column" form:"column"`
}

type CategoryDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type RoleDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UserDTO struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Email string  `json:"email"`
	Role  RoleDTO `json:"role"`
}

type Response struct {
	ID         int64         `json:"id"`
	Name       string        `json:"name"`
	IsActive   bool          `json:"is_active"`
	CreatedAt  string        `json:"created_at"`
	UpdatedAt  string        `json:"updated_at"`
	Categories []CategoryDTO `json:"categories"`
	Users      []UserDTO     `json:"users"`
}

func ToResponse(model company_model.Company) Response {
	var categories []CategoryDTO
	for _, c := range model.Categories {
		categories = append(categories, CategoryDTO{
			ID:   c.ID,
			Name: c.Name,
		})
	}

	var users []UserDTO
	for _, u := range model.CompanyUserRoles {
		users = append(users, UserDTO{
			ID:    u.UserID,
			Name:  u.User.Name,
			Email: u.User.Email,
			Role: RoleDTO{
				ID:   u.Role.ID,
				Name: u.Role.Name,
			},
		})
	}

	return Response{
		ID:         model.ID,
		Name:       model.Name,
		IsActive:   model.IsActive,
		CreatedAt:  helper.FormatDate(model.CreatedAt),
		UpdatedAt:  helper.FormatDate(model.UpdatedAt),
		Categories: categories,
		Users:      users,
	}
}
