package database

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//Function to connect to and return the database instance

func ConnectToDatabase() (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	dsn := os.Getenv("DSN_STRING")
	for retries := 5; retries > 0; retries-- {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db, nil
		}
		fmt.Printf("Failed to connect to database. Retrying in 5 seconds... (%d attempts left)\n", retries-1)
		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database after multiple attempts: %v", err)
}
