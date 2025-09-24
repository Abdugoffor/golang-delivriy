package seeder

import (
	"log"
	config "my-project/config"
	category_model "my-project/modul/category/model"
	"time"

	"gorm.io/gorm"
)

func CategorySeeder() {
	categories := []category_model.Category{
		{Name: "Elektronika", Slug: "elektronika", Order: 1},
		{Name: "Maishiy texnika", Slug: "maishiy-texnika", Order: 2},
		{Name: "Kiyim-kechak", Slug: "kiyim-kechak", Order: 3},
		{Name: "Kitoblar", Slug: "kitoblar", Order: 4},
		{Name: "O‘yinchoqlar", Slug: "oyinchoqlar", Order: 5},
		{Name: "Sport anjomlari", Slug: "sport-anjomlari", Order: 6},
		{Name: "Mebel", Slug: "mebel", Order: 7},
		{Name: "Idish-tovoq", Slug: "idish-tovoq", Order: 8},
		{Name: "Avto ehtiyot qismlar", Slug: "avto-ehtiyot-qismlar", Order: 9},
		{Name: "Kosmetika", Slug: "kosmetika", Order: 10},
	}

	for _, c := range categories {
		var exists category_model.Category
		if err := config.DB.Where("slug = ?", c.Slug).First(&exists).Error; err == gorm.ErrRecordNotFound {
			c.IsActive = true
			c.CreatedAt = time.Now()
			c.UpdatedAt = time.Now()

			if err := config.DB.Create(&c).Error; err != nil {
				log.Printf("❌ CategorySeeder insert error: %v", err)
			} else {
				log.Printf("✅ Category '%s' qo‘shildi", c.Name)
			}
		} else {
			log.Printf("ℹ️ Category '%s' allaqachon mavjud, o‘tkazib yuborildi", c.Name)
		}
	}
}
