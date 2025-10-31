package config

import (
	"fmt"
	"log"
	"os"

	"my-project/helper"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnect() *gorm.DB {
	driver := helper.ENV("DB_DRIVER")
	var db *gorm.DB
	var err error

	switch driver {
	case "postgres":
		// PostgreSQL DSN
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			helper.ENV("DB_HOST"),
			helper.ENV("DB_USER"),
			helper.ENV("DB_PASSWORD"),
			helper.ENV("DB_NAME"),
			helper.ENV("DB_PORT"),
			helper.ENV("DB_SSLMODE"),
			helper.ENV("DB_TIMEZONE"),
		)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("❌ Failed to connect to PostgreSQL:", err)
		}
		log.Println("✅ Connected to PostgreSQL")

	case "sqlite":
		// SQLite fayl yo‘lini olish
		sqlitePath := helper.ENV("DB_PATH")
		if sqlitePath == "" {
			sqlitePath = "data.db"
		}

		// Fayl mavjudligini tekshirish, agar yo‘q bo‘lsa yaratish
		if _, err := os.Stat(sqlitePath); os.IsNotExist(err) {
			file, createErr := os.Create(sqlitePath)
			if createErr != nil {
				log.Fatal("❌ Failed to create SQLite file:", createErr)
			}
			file.Close()
			log.Println("🆕 SQLite database file created:", sqlitePath)
		}

		db, err = gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})
		if err != nil {
			log.Fatal("❌ Failed to connect to SQLite:", err)
		}
		log.Println("✅ Connected to SQLite (" + sqlitePath + ")")

	default:
		log.Fatal("❌ Unknown DB_DRIVER. Please set DB_DRIVER=postgres or sqlite")
	}

	DB = db

	// Migration ishga tushadi
	RunMigrations()

	return db
}
