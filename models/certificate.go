package models
import("time")

// ---------------------
// Certificate
// ---------------------
type Certificate struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	User      *User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	CourseID  uint      `gorm:"not null" json:"course_id"`
	Course    *Course   `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"course,omitempty"`
	CertCode      string    `gorm:"size:100;unique;not null" json:"code"`
	URL       string    `gorm:"size:255" json:"url"`
	IssuedAt  time.Time `json:"issued_at"`
}