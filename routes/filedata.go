package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Swetabh333/trademarkia/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func GetUserFiles(db *gorm.DB, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("ID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}

		userIDUUID, ok := userID.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Invalid user ID format",
			})
			return
		}

		userIDStr := userIDUUID.String()
		cachedKey := fmt.Sprintf("user_files:%s", userIDStr)
		cachedFiles, err := redisClient.Get(c, cachedKey).Result()
		if err == nil {
			c.Data(http.StatusOK, "application/json", []byte(cachedFiles)) // data found in redis cache
			return
		}

		files := []models.Filedata{}
		fmt.Println(userIDUUID)
		if err := db.Where("user_id = ?", userIDUUID.String()).Find(&files).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to retreive files",
			})
			return
		}
		cachedData, err := json.Marshal(files)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to retreive files",
			})
			return
		}
		redisClient.Set(c, cachedKey, cachedData, 5*time.Minute)

		c.JSON(http.StatusOK, files)
	}
}

// Function for sending url of file
func Sharefile(db *gorm.DB, baseURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileID := c.Param("file_id")
		userID, ok := c.Get("ID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}

		file := models.Filedata{}

		if err := db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "File not found or you don't have permission"})
			return
		}

		publicURL := fmt.Sprintf("%s/%s", baseURL, file.Path)
		c.JSON(http.StatusOK, gin.H{"public_url": publicURL})
	}
}
