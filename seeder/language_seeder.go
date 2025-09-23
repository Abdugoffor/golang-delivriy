package seeder

import (
	"errors"
	"log"
	config "my-project/config"
	language_model "my-project/modul/language/model"

	"gorm.io/gorm"
)

func LanguageSeeder() {
	languages := []language_model.Language{
		{Name: "Uzbek", Slug: "uz", IsActive: true},
		{Name: "Russian", Slug: "ru", IsActive: true},
		{Name: "English", Slug: "en", IsActive: false},
	}

	for _, lang := range languages {
		var model language_model.Language

		err := config.DB.Where("slug = ?", lang.Slug).First(&model).Error

		if err == nil {
			continue // allaqachon mavjud
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("❌ DB error: %v\n", err)
			continue
		}

		if err := config.DB.Create(&lang).Error; err != nil {
			log.Printf("❌ Failed to insert language %s: %v\n", lang.Slug, err)
		}
	}
	log.Println("✅ Languages seeded successfully")
}
