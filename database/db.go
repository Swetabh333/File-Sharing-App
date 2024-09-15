package databse

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

//Function to connect to and return the database instance

func ConnectToDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("DSN_STRING")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}
