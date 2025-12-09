package models

import (
	"time"
)

type Payment struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint      `gorm:"index;not null" json:"user_id"`
	User          User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	CourseID      uint      `gorm:"index;not null" json:"course_id"`
	Course        Course    `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"course"`
	Amount        float64   `gorm:"not null" json:"amount"`
	Status        string    `gorm:"type:enum('pending','success','failed');default:'pending'" json:"status"`
	TransactionID string    `gorm:"size:100;unique" json:"transaction_id"`
	PaymentMethod string    `gorm:"size:50" json:"payment_method"` // e.g., card, upi
	CreatedAt     time.Time `json:"created_at"`
    DiscountApplied float64 `json:"discount_applied"` // Discount applied if any
}
