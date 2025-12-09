package models

import "time"

// Assignment represents a test/quiz for a course
type Assignment struct {
    ID          uint          `gorm:"primaryKey;autoIncrement" json:"id"`
    Title       string        `gorm:"size:255;not null" json:"title"`
    Description string        `gorm:"type:text" json:"description"`

    CourseID uint    `gorm:"not null" json:"course_id"`
    Course   Course  `gorm:"foreignKey:CourseID"`

    TeacherID uint           `gorm:"not null" json:"teacher_id"` // Foreign key field
    Teacher   TeacherProfile `gorm:"foreignKey:TeacherID"`       // Explicit foreign key link

    Questions []Question `gorm:"constraint:OnDelete:CASCADE" json:"questions"`

    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}


// Question represents an MCQ inside an assignment
type Question struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AssignmentID  uint      `gorm:"not null" json:"assignment_id"`
	Assignment   Assignment `gorm:"foreignKey:AssignmentID"`

	Text  string    `gorm:"type:text;not null" json:"question_text"`
	Options       []Option  `gorm:"constraint:OnDelete:CASCADE" json:"options"`
	CorrectOption uint      `json:"correct_option"` // Reference to Option.ID
}

// Option represents a single choice for a Question
type Option struct {
	ID         uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	QuestionID uint     `gorm:"not null" json:"question_id"`
	Question   Question `gorm:"foreignKey:QuestionID"`
	IsCorrect  bool   `gorm:"default:false" json:"is_correct"`
	Text       string   `gorm:"type:text;not null" json:"text"`
}

// AssignmentSubmission represents student's attempt
type Submission struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AssignmentID uint      `gorm:"not null" json:"assignment_id"`
	Assignment   Assignment `gorm:"foreignKey:AssignmentID"`

	UserID uint `gorm:"not null" json:"user_id"`
	User   User `gorm:"foreignKey:UserID"`

	Score      int       `json:"score"`       // Calculated score
	SubmittedAt time.Time `json:"submitted_at"`
}
