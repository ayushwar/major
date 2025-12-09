package models

import("time")

type Department struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"size:255;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	// Thumbnail / Image URL for UI
	ThumbnailURL string `gorm:"size:255" json:"thumbnail_url"`

	// Relations (1 Department → Many Courses)
	Courses []Course `gorm:"foreignKey:DepartmentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"courses,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
