package seeder

import (
	"my-project/config"
	user_model "my-project/modul/user/model"
)

func CompanyUserRoleSeeder() {
	relations := []user_model.CompanyUserRole{
		{CompanyID: 1, UserID: 1, RoleID: 1, IsActive: true}, // Admin
		{CompanyID: 1, UserID: 2, RoleID: 2, IsActive: true}, // Seller
		{CompanyID: 2, UserID: 3, RoleID: 1, IsActive: true}, // Buyer
		{CompanyID: 2, UserID: 5, RoleID: 1, IsActive: true}, // Buyer
		{CompanyID: 3, UserID: 4, RoleID: 2, IsActive: true}, // Buyer
		{CompanyID: 3, UserID: 2, RoleID: 1, IsActive: true}, // Buyer
		{CompanyID: 3, UserID: 6, RoleID: 2, IsActive: true}, // Buyer
		{CompanyID: 4, UserID: 5, RoleID: 1, IsActive: true}, // Buyer
		{CompanyID: 4, UserID: 1, RoleID: 2, IsActive: true}, // Buyer
		{CompanyID: 5, UserID: 6, RoleID: 2, IsActive: true}, // Buyer
		{CompanyID: 5, UserID: 3, RoleID: 1, IsActive: true}, // Buyer
		{CompanyID: 5, UserID: 4, RoleID: 2, IsActive: true}, // Buyer
	}

	for _, r := range relations {
		config.DB.FirstOrCreate(&r, user_model.CompanyUserRole{
			CompanyID: r.CompanyID,
			UserID:    r.UserID,
			RoleID:    r.RoleID,
		})
	}
}
