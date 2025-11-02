package config

import (
	"log"
	category_model "my-project/modul/category/model"
	product_model "my-project/modul/product/model"
)

func RunMigrations() {
	models := []interface{}{
		&product_model.Product{},
		&category_model.Category{},
	}
	err := DB.AutoMigrate(models...)

	if err != nil {
		log.Fatal("❌ Failed to run migrations: ", err)
	}

	log.Println("✅ Migrations completed")

	// CreateHistoryTriggers(DB, models)
}
