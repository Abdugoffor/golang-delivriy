package config

import (
	"fmt"
	"log"
	"my-project/helper"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnect() *gorm.DB {
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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("❌ Failed to connect to database: ", err)
	}

	DB = db
	log.Println("✅ PostgreSQL connected")

	RunMigrations()

	return db
}
