package models

import (
	"time"
)

type CollegeVerification struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	College    string    `gorm:"size:100;not null" json:"college"`
	Document   string    `gorm:"size:255;not null" json:"document"` // uploaded document
	Status     string    `gorm:"type:enum('pending','verified','rejected');default:'pending'" json:"status"`
	VerifiedAt time.Time `json:"verified_at"`
	CreatedAt  time.Time `json:"created_at"`
}
