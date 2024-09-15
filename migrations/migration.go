package main

import (
	databse "github.com/Swetabh333/trademarkia/database"
	models "github.com/Swetabh333/trademarkia/models"
	"github.com/joho/godotenv"
	"log"
)

//This function runs the migrations before the project is built and run so
//that all the required tables and relations are created in our postgres databse

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	db, err := databse.ConnectToDatabase()
	if err != nil {
		log.Fatalf("Error connectng to Database: %s", err)
	}
	err = db.AutoMigrate(&models.User{}, &models.Filedata{})

	if err != nil {
		log.Fatalf("Failed to migrate")
	}
	log.Println("Migration completed.")
}
