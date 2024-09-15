package routes

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Swetabh333/trademarkia/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	uploadDir   = "./uploads" //local directory for storing files
	maxFileSize = 10 << 20    // 10MB
)

//Handler function for handling file uploads

func UploadHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDinterface, ok := c.Get("ID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User not authenticated",
			})
			return
		}
		userID, ok := userIDinterface.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
			return
		}
		if err := c.Request.ParseMultipartForm(maxFileSize); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "File size should be less than 10MB",
			})
			return
		}
		files := c.Request.MultipartForm.File["files"]
		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "No files uploaded",
			})
			return
		}

		wg := sync.WaitGroup{}
		results := make(chan *models.Filedata, len(files))
		errors := make(chan error, len(files))

		for _, file := range files {
			wg.Add(1)
			go func(file *multipart.FileHeader) {
				defer wg.Done()
				result, err := processFile(db, file, userID)
				if err != nil {
					errors <- err
				} else {
					results <- result
				}
			}(file)
		}

		go func() {
			wg.Wait()
			close(results)
			close(errors)
		}()
		uploadedFiles := []*models.Filedata{}
		errorMessages := []string{}

		for i := 0; i < len(files); i++ {
			select {
			case result := <-results:
				uploadedFiles = append(uploadedFiles, result)

			case err := <-errors:
				errorMessages = append(errorMessages, err.Error())
			}
		}
		response := gin.H{
			"message": fmt.Sprintf("%d files uploaded successfully", len(uploadedFiles)),
			"files":   uploadedFiles,
		}
		if len(errorMessages) > 0 {
			response["errors"] = errorMessages
		}

		c.JSON(http.StatusOK, response)
	}
}

//Function to create a new file name , make sure the directory for storing exists and for storing the file metadata

func processFile(db *gorm.DB, file *multipart.FileHeader, userID uuid.UUID) (*models.Filedata, error) {

	filename := uuid.New().String() + filepath.Ext(file.Filename)
	filepath := filepath.Join(uploadDir, filename)
	fmt.Printf("filename: %s , filepath: %s \n", filename, filepath)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to save file %w", err)
	}

	if err := saveFile(file, filepath); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	filedata := &models.Filedata{
		ID:         uuid.New(),
		UserID:     userID,
		Name:       file.Filename,
		Path:       filepath,
		Size:       file.Size,
		UploadedAt: time.Now(),
	}
	if err := db.Create(filedata).Error; err != nil {
		// If database insert fails, delete the uploaded file
		os.Remove(filepath)
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}
	return filedata, nil
}

//Function to copy uploaded source file to our upload location in destination

func saveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
