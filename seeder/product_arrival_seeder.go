package seeder

import (
	"log"
	"math/rand"
	config "my-project/config"
	product_model "my-project/modul/product/model"
	product_arrival_model "my-project/modul/product_arrival/model"
	"time"
)

// Exceldan keladigan namunaviy ma'lumotlar (Author, ProductName, Count, Price)
type ArrivalInput struct {
	Author      string
	ProductName string
	Count       int
	Price       int
}

func ProductArrivalSeeder() {
	var products []product_model.Product
	if err := config.DB.Find(&products).Error; err != nil {
		log.Printf("❌ Productlarni olishda xatolik: %v", err)
		return
	}

	for _, product := range products {
		// Tasodifiy miqdor va narx
		count := rand.Intn(50) + 1 // 1–50 dona
		price := product.Price     // mavjud narxi
		sum := count * price

		// Oldingi qoldiqni topish
		var lastLog product_arrival_model.ProductArrivalLog
		config.DB.Where("product_id = ?", product.ID).Order("id desc").First(&lastLog)
		quantityBefore := lastLog.QuantityAfter

		// 🔹 ProductArrival yozuvi
		arrival := product_arrival_model.ProductArrival{
			Author:    "SeederBot",
			ProductID: product.ID,
			Count:     count,
			Sum:       sum,
			CreatedAt: time.Now(),
		}
		if err := config.DB.Create(&arrival).Error; err != nil {
			log.Printf("❌ ProductArrival error: %v", err)
			continue
		}

		// 🔹 ProductArrivalLog yozuvi
		logEntry := product_arrival_model.ProductArrivalLog{
			Type:           "arrival",
			ProductID:      product.ID,
			QuantityBefore: quantityBefore,
			Quantity:       count,
			QuantityAfter:  quantityBefore + count,
			Sum:            sum,
			CreatedAt:      time.Now(),
		}
		if err := config.DB.Create(&logEntry).Error; err != nil {
			log.Printf("❌ ProductArrivalLog error: %v", err)
			continue
		}

		log.Printf("📦 Prihod: %s (%d dona, %d so‘m)", product.Name, count, sum)
	}
}
