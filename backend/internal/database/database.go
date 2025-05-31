package database

import (
	"log"
	"stockify/internal/core"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connection successful!")

	if err := db.Exec("SELECT 1").Error; err != nil {
		log.Fatal("Database ping failed:", err)
	}
	log.Println("Database ping successful!")

	log.Println("Running database migrations...")
	if err := db.AutoMigrate(&core.Stock{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrations successful!")
	return db
}
