package models

import (
	"time"
)

// ---------------------
// Course
// ---------------------


type Course struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	// Teacher relationship: TeacherID stores the User ID (from JWT)
	TeacherID uint `gorm:"not null;index:idx_teacher_courses" json:"teacher_id"`
	
	// ✅ CRITICAL FIX: Explicitly reference the UserID column in the profile table
	Teacher *TeacherProfile `gorm:"foreignKey:TeacherID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"teacher,omitempty"`

	// Department Link (Required for Course)
	DepartmentID uint  `gorm:"not null;index:idx_dept_courses" json:"department_id" binding:"required"`
	Department *Department `gorm:"foreignKey:DepartmentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"department,omitempty"`

	// Course fields (Ensure these match your input fields)
	Title string `gorm:"size:255;not null" json:"title" binding:"required"`
	Code  string `gorm:"size:50;uniqueIndex;not null" json:"code" binding:"required"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	Credits int `gorm:"not null;default:3" json:"credits" binding:"required,min=1,max=10"`
	
	// Add other fields you need here (e.g., Level, Language)
	// Example: Level string `gorm:"size:50" json:"level"` 

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ... Lecture struct follows ...


// ---------------------
// Lecture (stored on YouTube)
// ---------------------
// LectureStatus defines the possible states of a lecture
type LectureStatus string

const (
	LectureStatusProcessing LectureStatus = "processing"
	LectureStatusReady  LectureStatus = "ready"
	LectureStatusFailed LectureStatus = "failed"
	LectureStatusUploading LectureStatus = "uploading"
)

type Lecture struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	// Foreign Key linking to the Course
	CourseID uint `gorm:"not null;index:idx_course_lectures" json:"course_id"`
	Course*Course `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"course,omitempty"`

	// Metadata provided by the teacher
	Title string `gorm:"size:255;not null" json:"title" binding:"required"`
	Description string `gorm:"type:text" json:"description"`
	OrderIndex int `gorm:"default:0" json:"order_index"` // For sequencing lectures

	// YouTube Data Storage
	YouTubeVideoID string `gorm:"size:50;uniqueIndex" json:"youtube_video_id,omitempty"` // e.g., "dQw4w9WgXcQ"
	YouTubeURL string `gorm:"size:255" json:"youtube_url,omitempty"` // Full URL
	Duration string `gorm:"size:50" json:"duration,omitempty"` // ISO 8601 duration (e.g., "PT15M33S")
	ThumbnailURL string `gorm:"size:255" json:"thumbnail_url,omitempty"` // YouTube thumbnail

	// Status tracking
	Status LectureStatus `gorm:"size:20;default:'uploading'" json:"status"`
	ErrorMessage string `gorm:"type:text" json:"error_message,omitempty"` // Store error details if upload fails

	// Metadata
	FileSize int64 `json:"file_size,omitempty"`// Size in bytes
	MimeType string `gorm:"size:100" json:"mime_type,omitempty"`
	UploadedBy uint`gorm:"not null" json:"uploaded_by"` // User ID who uploaded
	ViewCount int64`gorm:"default:0" json:"view_count"`
	IsPublished bool`gorm:"default:false" json:"is_published"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Lecture) TableName() string {
	return "lectures"
}