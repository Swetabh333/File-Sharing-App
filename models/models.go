package models

import (
	"time"

	"github.com/google/uuid"
)

// Model for storing user data in the database
type User struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Name      string    `gorm:"size:255;not null;unique"`
	Email     string    `gorm:"size:255;unique"`
	Password  string    `gorm:"size:255;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Model for storing file metadata in the database
type Filedata struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	UserID     uuid.UUID `gorm:"not null"`
	Name       string    `gorm:"size:255;not null"`
	Path       string    `gorm:"size:255;not null"`
	Size       int64     `gorm:"size255;not null"`
	UploadedAt time.Time
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	User       User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
