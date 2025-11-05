package config

import (
	"log"
	category_model "my-project/modul/category/model"
	company_model "my-project/modul/company/model"
	product_model "my-project/modul/product/model"
	user_model "my-project/modul/user/model"
)

func RunMigrations() {
	models := []interface{}{
		&user_model.Role{},            // 1️⃣ avval rollar
		&user_model.User{},            // 2️⃣ userlar
		&company_model.Company{},      // 3️⃣ company
		&category_model.Category{},    // 4️⃣ category
		&user_model.CompanyUserRole{}, // 5️⃣ many-to-many bog‘lanma

		&product_model.Product{},
	}
	err := DB.AutoMigrate(models...)

	if err != nil {
		log.Fatal("❌ Failed to run migrations: ", err)
	}

	log.Println("✅ Migrations completed")

	// CreateHistoryTriggers(DB, models)
}
