package seeder

import (
	"log"
	"math/rand"
	config "my-project/config"
	"my-project/helper"
	product_model "my-project/modul/product/model"
	"time"
)

var productNames = []string{
	"Snickers", "Twix", "Mars", "Bounty", "KitKat",
	"Choco Pie", "Nutella", "Milka", "Toblerone", "Ferrero Rocher",
	"Pepsi", "Coca Cola", "Fanta", "Sprite", "7Up",
	"Mountain Dew", "Red Bull", "Monster Energy", "Lipton Ice Tea", "Nestea",
	"Lays", "Pringles", "Doritos", "Cheetos", "Ruffles",
	"Colgate", "Sensodyne", "CloseUp", "Aquafresh", "Oral-B",
	"Dove Soap", "Lux Soap", "Lifebuoy", "Dettol", "Palmolive",
	"Nescafe", "Jacobs", "Lavazza", "Tchibo", "Starbucks",
	"Ahmad Tea", "Lipton Tea", "Greenfield Tea", "Akbar Tea", "Curtis Tea",
	"Head & Shoulders", "Pantene", "Sunsilk", "Garnier", "L'Oreal",
	"Adidas Shoes", "Nike Shoes", "Puma Shoes", "Reebok Shoes", "New Balance Shoes",
	"Samsung Galaxy S21", "iPhone 13", "Xiaomi Redmi Note 10", "OnePlus 9", "Huawei P40",
	"Dell Laptop", "HP Laptop", "Lenovo Laptop", "Asus Laptop", "Acer Laptop",
	"Canon Camera", "Nikon Camera", "Sony Camera", "Fujifilm Camera", "GoPro",
	"LG TV", "Samsung TV", "Sony Bravia TV", "Panasonic TV", "Philips TV",
	"PS5 Console", "Xbox Series X", "Nintendo Switch", "Logitech Mouse", "Razer Keyboard",
	"Bosch Fridge", "LG Fridge", "Samsung Fridge", "Whirlpool Fridge", "Haier Fridge",
	"Bosch Washing Machine", "Samsung Washing Machine", "LG Washing Machine", "Indesit Washing Machine", "Candy Washing Machine",
	"Toyota Corolla", "Honda Civic", "Hyundai Elantra", "Kia Sportage", "Chevrolet Spark",
	"BMW X5", "Mercedes C-Class", "Audi A6", "Volkswagen Golf", "Tesla Model 3",
	"Adidas Shoes", "Nike Shoes", "Puma Shoes", "Reebok Shoes", "New Balance Shoes",
	"Samsung Galaxy S21", "iPhone 13", "Xiaomi Redmi Note 10", "OnePlus 9", "Huawei P40",
	"Dell Laptop", "HP Laptop", "Lenovo Laptop", "Asus Laptop", "Acer Laptop",
	"Canon Camera", "Nikon Camera", "Sony Camera", "Fujifilm Camera", "GoPro",
	"LG TV", "Samsung TV", "Sony Bravia TV", "Panasonic TV", "Philips TV",
	"PS5 Console", "Xbox Series X", "Nintendo Switch", "Logitech Mouse", "Razer Keyboard",
	"Bosch Fridge", "LG Fridge", "Samsung Fridge", "Whirlpool Fridge", "Haier Fridge",
	"Bosch Washing Machine", "Samsung Washing Machine", "LG Washing Machine", "Indesit Washing Machine", "Candy Washing Machine",
	"Toyota Corolla", "Honda Civic", "Hyundai Elantra", "Kia Sportage", "Chevrolet Spark",
	"BMW X5", "Mercedes C-Class", "Audi A6", "Volkswagen Golf", "Tesla Model 3",
}

func ProductSeeder() {
	for i, name := range productNames {
		slug := helper.Slug(name)

		product := product_model.Product{
			Name:      name,
			Slug:      slug,
			Price:     rand.Int63n(300),
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := config.DB.Create(&product).Error; err != nil {
			log.Printf("❌ ProductSeeder insert error: %v", err)
		} else {
			log.Printf("✅ (%d) Product qo‘shildi: %s", i+1, product.Name)
		}
	}
}
