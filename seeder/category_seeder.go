package seeder

import (
	"my-project/config"
	category_model "my-project/modul/category/model"
)

func CategorySeeder() {
	categories := []category_model.Category{
		{Name: "Electronics", Slug: "electronics", CompanyID: 1, IsActive: true},
		{Name: "Clothing", Slug: "clothing", CompanyID: 2, IsActive: true},
		{Name: "Home", Slug: "home", CompanyID: 6, IsActive: true},
		{Name: "Beauty", Slug: "beauty", CompanyID: 3, IsActive: true},
		{Name: "Sports", Slug: "sports", CompanyID: 4, IsActive: true},
		{Name: "Toys", Slug: "toys", CompanyID: 5, IsActive: true},
	}

	for _, cat := range categories {
		config.DB.FirstOrCreate(&cat, category_model.Category{Slug: cat.Slug, CompanyID: cat.CompanyID})
	}
}
