package config

import (
	"log"
	product_model "my-project/modul/product/model"
)

func RunMigrations() {
	models := []interface{}{
		&product_model.Product{},
	}
	err := DB.AutoMigrate(models...)

	if err != nil {
		log.Fatal("❌ Failed to run migrations: ", err)
	}

	log.Println("✅ Migrations completed")

	// CreateHistoryTriggers(DB, models)
}
