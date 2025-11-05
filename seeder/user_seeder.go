package seeder

import (
	"my-project/config"
	user_model "my-project/modul/user/model"
)

func UserSeeder() {
	users := []user_model.User{
		{Name: "Admin", Email: "admin@gmail.com", Password: "123456789", IsActive: true},
		{Name: "User", Email: "user@gmail.com", Password: "123456789", IsActive: true},
		{Name: "Guest", Email: "guest@gmail.com", Password: "123456789", IsActive: true},
		{Name: "Seller", Email: "seller@gmail.com", Password: "123456789", IsActive: true},
		{Name: "Buyer", Email: "buyer@gmail.com", Password: "123456789", IsActive: true},
		{Name: "Manager", Email: "manager@gmail.com", Password: "123456789", IsActive: true},
	}

	for _, user := range users {
		config.DB.FirstOrCreate(&user, user_model.User{Email: user.Email})
	}
}
