package seeder

import (
	"my-project/config"
	user_model "my-project/modul/user/model"
)

func RoleSeeder() {
	roles := []user_model.Role{
		{Name: "admin", IsActive: true},
		{Name: "seller", IsActive: true},
		{Name: "buyer", IsActive: true},
	}

	for _, role := range roles {
		config.DB.FirstOrCreate(&role, user_model.Role{Name: role.Name})
	}
}
