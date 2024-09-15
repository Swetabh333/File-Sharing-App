package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Swetabh333/trademarkia/database"
	"github.com/Swetabh333/trademarkia/middleware"
	"github.com/Swetabh333/trademarkia/models"
	"github.com/Swetabh333/trademarkia/routes"
	"github.com/Swetabh333/trademarkia/utils"

	"github.com/gin-gonic/gin"

	//	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

const (
	uploadDir = "./uploads"
	baseURL   = "http://localhost:8080"
)

// background deletion
func initFileDeleteWorker(ctx context.Context, db *gorm.DB) {
	worker := utils.NewFileDeleteWorker(
		db,
		1*time.Hour,     // Run every hour
		30*24*time.Hour, // Files expire after 30 days
	)

	go worker.Start(ctx)
}

func main() {
	//Loading the environament variables
	//	err := godotenv.Load(".env")
	//	if err != nil {
	//		log.Fatalf("Error loading .env file: %s", err)
	//	}

	//connecting to database
	db, err := database.ConnectToDatabase()
	if err != nil {
		log.Fatalf("Could not connect to the database")
	}
	fmt.Println("Connected to db")

	if err != nil {
		log.Fatalf("Error connectng to Database: %s", err)
	}
	err = db.AutoMigrate(&models.User{}, &models.Filedata{})

	if err != nil {
		log.Fatalf("Failed to migrate : %s", err)
	}
	log.Println("Migration completed.")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize the file deletion worker
	initFileDeleteWorker(ctx, db)

	redis_client := database.ConnectToRedis()
	//Setting up our gin http server
	router := gin.Default()
	router.POST("/register", routes.HandleRegistration(db))
	router.POST("/login", routes.HandleLogin(db))
	router.GET("/verify", middleware.CheckAuthentication, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "Verified",
		})
	})
	router.Static("/uploads", uploadDir)
	router.POST("/upload", middleware.CheckAuthentication, routes.UploadHandler(db))
	router.GET("/files", middleware.CheckAuthentication, routes.GetUserFiles(db, redis_client))
	router.GET("/share/:file_id", middleware.CheckAuthentication, routes.Sharefile(db, baseURL))
	router.POST("/search", middleware.CheckAuthentication, routes.SearchFiles(db))
	router.Run("0.0.0.0:8080")
}
