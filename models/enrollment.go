package models

import (
	"time"
)

type Enrollment struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint       `gorm:"not null;uniqueIndex:idx_user_course" json:"user_id"`
	User         *User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	CourseID     uint       `gorm:"not null;uniqueIndex:idx_user_course" json:"course_id"`
	Course       *Course    `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"course,omitempty"`

	Progress     float32    `gorm:"default:0" json:"progress"` // % completed
	Status       string     `gorm:"size:20;default:'active'" json:"status"`
	EnrolledAt   time.Time  `json:"enrolled_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	CertificateID *string   `gorm:"size:100" json:"certificate_id,omitempty"`

	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
