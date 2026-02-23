package	models

import (
	"time"

	"gorm.io/gorm"
)

// Video represents a YouTube video linked to a batch
type Video struct {
	ID             uint           `gorm:"primaryKey;autoIncrement"`
	Title          string         `gorm:"type:text;not null"`
	Description    string         `gorm:"type:text"`
	YouTubeVideoID string         `gorm:"type:varchar(50);not null"`
	BatchID        uint           `gorm:"not null;index"`
	TeacherID      uint           `gorm:"not null;index"`
	CreatedAt      time.Time      `gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}