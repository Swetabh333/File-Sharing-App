package routes

import (
	"time"

	"github.com/Swetabh333/trademarkia/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SearchParams struct {
	Name       string    `form:"name"`
	UploadDate time.Time `form:"upload_date" time_format:"2006-01-02"`
	Page       int       `form:"page,default=1"`
	PageSize   int       `form:"page_size,default=10"`
}

func SearchFiles(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params SearchParams
		if err := c.ShouldBindQuery(&params); err != nil {
			c.JSON(400, gin.H{"error": "Invalid search parameters"})
			return
		}

		userID, _ := c.Get("ID")

		query := db.Model(&models.Filedata{}).Where("user_id = ?", userID)

		// Apply search filters
		if params.Name != "" {
			query = query.Where("name LIKE ?", "%"+params.Name+"%")
		}
		if !params.UploadDate.IsZero() {
			query = query.Where("DATE(uploaded_at) = DATE(?)", params.UploadDate)
		}

		// Count total results
		var total int64
		query.Count(&total)

		// Apply pagination
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)

		// Execute query
		var files []models.Filedata
		if err := query.Find(&files).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to search files"})
			return
		}

		c.JSON(200, gin.H{
			"total":     total,
			"page":      params.Page,
			"page_size": params.PageSize,
			"files":     files,
		})
	}
}
