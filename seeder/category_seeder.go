package seeder

import (
	"my-project/config"
	category_model "my-project/modul/category/model"
	"strings"
)

func CategorySeeder() {
	names := []string{
		"Electronics", "Clothing", "Home", "Beauty", "Sports", "Toys",
		"Books", "Shoes", "Bags", "Accessories",
		"Furniture", "Kitchen", "Garden", "Automotive", "Health",
		"Food", "Drinks", "Office", "Pets", "Kids",
		"Gadgets", "Tools", "Hardware", "Software", "Music",
		"Movies", "Games", "Travel", "Outdoor", "Fitness",
		"Jewelry", "Watches", "Phones", "Laptops", "Tablets",
		"TV", "Cameras", "Printers", "Networking", "Security",
		"Lighting", "Decor", "Art", "Craft", "Stationery",
		"Cleaning", "Laundry", "Storage", "Baby", "Maternity",
		"Medical", "Supplements", "Cosmetics", "Fragrance", "Haircare",
		"Skincare", "Makeup", "Nails", "Barber", "Salon",
		"Gym", "Yoga", "Cycling", "Running", "Hiking",
		"Camping", "Fishing", "Hunting", "Swimming", "Diving",
		"Winter", "Summer", "Rain", "Smart Home", "AI Devices",
		"Drones", "Robotics", "VR", "AR", "Gaming Gear",
		"Streaming", "Cloud", "Hosting", "Domains", "SEO",
		"Marketing", "Finance", "Crypto", "Banking", "Insurance",
		"Education", "Courses", "Languages", "Consulting", "Services",
	}

	var categories []category_model.Category
	for p := 0; p < 2500; p++ {

		for i, name := range names {
			categories = append(categories, category_model.Category{
				Name:      name,
				Slug:      strings.ToLower(strings.ReplaceAll(name, " ", "-")),
				CompanyID: int64((i % 5) + 1),
				IsActive:  true,
			})
		}

		config.DB.Create(&categories)
	}
}
