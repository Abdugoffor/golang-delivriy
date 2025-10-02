package config

import (
	"log"
	category_model "my-project/modul/category/model"
	language_model "my-project/modul/language/model"
	product_model "my-project/modul/product/model"
	product_arrival_model "my-project/modul/product_arrival/model"
)

func RunMigrations() {
	err := DB.AutoMigrate(
		&language_model.Language{},
		&category_model.Category{},
		&product_model.Product{},
		&product_arrival_model.ProductArrival{},
		&product_arrival_model.ProductArrivalLog{},
	)

	if err != nil {
		log.Fatal("❌ Failed to run migrations: ", err)
	}

	log.Println("✅ Migrations completed")

	CreateHistoryTriggers([]string{"categories", "products", "product_arrivals", "languages"})
}
