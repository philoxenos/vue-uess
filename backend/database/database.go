package database

import (
	"log"
	"mis-system/models"

	"github.com/glebarez/sqlite" // Pure Go SQLite driver
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase initializes the database connection
func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("mis.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate the database
	err = database.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.AuthAudit{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	DB = database
	log.Println("Database connected successfully")
}
