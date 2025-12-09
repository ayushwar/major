package models
import("time")


type Progress struct {
    ID        uint    `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID    uint    `gorm:"not null" json:"user_id"`
    CourseID  uint    `gorm:"not null" json:"course_id"`
    CompletedAssignments int `json:"completed_assignments"`
    TotalAssignments     int `json:"total_assignments"`
    PercentageComplete   float32 `json:"percentage_complete"`

    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
