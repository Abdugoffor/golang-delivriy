package config

import (
	"log"
	category_model "my-project/modul/category/model"
	product_model "my-project/modul/product/model"
	user_model "my-project/modul/user/model"
)

func RunMigrations() {
	models := []interface{}{
		user_model.User{},
		category_model.Category{},
		product_model.Product{},
	}

	err := DB.AutoMigrate(models...)
	{
		if err != nil {
			log.Fatal("❌ Failed to run migrations: ", err)
		}
	}

	log.Println("✅ Migrations completed")

	CreateHistoryTriggers(DB, models)
}
