package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Swetabh333/trademarkia/models"
	"gorm.io/gorm"
)

type FileDeleteWorker struct {
	db           *gorm.DB
	interval     time.Duration
	expiryPeriod time.Duration
}

// Creates a new delete worker

func NewFileDeleteWorker(db *gorm.DB, interval time.Duration, expiryPeriod time.Duration) *FileDeleteWorker {
	return &FileDeleteWorker{
		db:           db,
		interval:     interval,
		expiryPeriod: expiryPeriod,
	}
}

func (w *FileDeleteWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("File deletion worker stopped")
			return
		case <-ticker.C:
			w.deleteExpiredFiles()
		}
	}
}

func (w *FileDeleteWorker) deleteExpiredFiles() {
	var expiredFiles []models.Filedata
	expiryDate := time.Now().Add(-w.expiryPeriod)

	if err := w.db.Where("uploaded_at < ?", expiryDate).Find(&expiredFiles).Error; err != nil {
		log.Printf("Error fetching expired files: %v", err)
		return
	}

	for _, file := range expiredFiles {
		// Delete from local storage
		filePath := fmt.Sprintf("./%s", file.Path)
		err := os.Remove(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("File %s already deleted from storage", filePath)
			} else {
				log.Printf("Error deleting file %s from storage: %v", filePath, err)
				continue
			}
		}

		// Delete from database
		if err := w.db.Delete(&file).Error; err != nil {
			log.Printf("Error deleting file metadata for %s: %v", file.Path, err)
		} else {
			log.Printf("Successfully deleted expired file: %s", file.Path)
		}
	}
}
