package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Swetabh333/trademarkia/database"
	"github.com/Swetabh333/trademarkia/middleware"
	"github.com/Swetabh333/trademarkia/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	//Loading the environament variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	//connecting to database
	db, err := databse.ConnectToDatabase()
	if err != nil {
		log.Fatalf("Could not connect to the database")
	}
	fmt.Println("Connected to db")

	//Setting up our gin http server
	router := gin.Default()
	router.POST("/register", routes.HandleRegistration(db))
	router.POST("/login", routes.HandleLogin(db))
	router.GET("/verify", middleware.CheckAuthentication, func(c *gin.Context) {
		id, _ := c.Get("ID")
		c.JSON(http.StatusOK, gin.H{
			"msg": fmt.Sprintf("Type of id is %T,", id),
		})
	})
	router.POST("/upload", middleware.CheckAuthentication, routes.UploadHandler(db))
	router.Run("localhost:8080")
}
