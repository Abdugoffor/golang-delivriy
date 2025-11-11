package config

import (
	"fmt"
	"log"
	"os"

	"my-project/helper"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnect() *gorm.DB {
	driver := helper.ENV("DB_DRIVER")
	var db *gorm.DB
	var err error

	switch driver {
	// 🟩 PostgreSQL
	case "postgres":
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
		{
			if err != nil {
				log.Fatal("❌ Failed to connect to PostgreSQL:", err)
			}
		}

		log.Println("✅ Connected to PostgreSQL")

	// 🟨 MySQL
	case "mysql":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			helper.ENV("DB_USER"),
			helper.ENV("DB_PASSWORD"),
			helper.ENV("DB_HOST"),
			helper.ENV("DB_PORT"),
			helper.ENV("DB_NAME"),
		)

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		{
			if err != nil {
				log.Fatal("❌ Failed to connect to MySQL:", err)
			}
		}

		log.Println("✅ Connected to MySQL")

	// 🟦 SQLite
	case "sqlite":
		sqlitePath := helper.ENV("DB_PATH")
		{
			if sqlitePath == "" {
				sqlitePath = "data.db"
			}
		}

		// Faylni tekshirish yoki yaratish
		if _, err := os.Stat(sqlitePath); os.IsNotExist(err) {
			file, createErr := os.Create(sqlitePath)

			if createErr != nil {
				log.Fatal("❌ Failed to create SQLite file:", createErr)
			}

			file.Close()

			log.Println("🆕 SQLite file created:", sqlitePath)
		}

		db, err = gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})
		{
			if err != nil {
				log.Fatal("❌ Failed to connect to SQLite:", err)
			}
		}

		log.Println("✅ Connected to SQLite (" + sqlitePath + ")")

	default:
		log.Fatal("❌ Unknown DB_DRIVER. Please set DB_DRIVER=postgres, mysql or sqlite")
	}

	DB = db

	RunMigrations()

	return db
}
