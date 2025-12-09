package models

import (
    "time"
)

// ---------------------
// User (auth & role)
// ---------------------
// Master entity: sab users (student, teacher, admin)
// ka record yahi rakha jata hai.
type User struct {
    ID          uint          `gorm:"primaryKey;autoIncrement" json:"id"`
    Name        string        `gorm:"size:100;not null" json:"name"`
    Email       string        `gorm:"size:100;unique;not null" json:"email"`
    Password string `gorm:"size:255;not null" json:"password,omitempty"`
    Role        string        `gorm:"type:enum('student','teacher','admin');default:'student'" json:"role"`

    // Auth / verification
    IsVerified  bool          `gorm:"default:false" json:"is_verified"`
    ResetToken  string        `gorm:"size:255" json:"-"`
    ResetExpiry time.Time     `json:"-"`

    // One-to-one relations
    Profile         *Profile        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"profile,omitempty"`
    TeacherProfile  *TeacherProfile `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"teacher_profile,omitempty"`
    
    // One-to-many
    Enrollments  []Enrollment   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"enrollments,omitempty"`
    Payments     []Payment      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"payments,omitempty"`
    Certificates []Certificate  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"certificates,omitempty"`

    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// ---------------------
// Student / General Profile
// ---------------------
type Profile struct {
    ID      uint    `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID  uint    `gorm:"uniqueIndex;not null" json:"user_id"`
    User    *User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`

    Image   string  `gorm:"size:255" json:"image"`
    College string  `gorm:"size:100;not null" json:"college"`    // mandatory for students
    Bio     string  `gorm:"size:255" json:"bio,omitempty"`
    
    Verified   bool   `gorm:"default:false" json:"verified"`        // verified by admin
    StudentID  string `gorm:"size:20;unique;not null" json:"student_id"` // e.g., 0101NT2501
    
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// ---------------------
// Teacher Profile
// ---------------------
// type TeacherProfile struct {
//     ID      uint    `gorm:"primaryKey;autoIncrement" json:"id"`
//     UserID  uint    `gorm:"uniqueIndex;not null" json:"user_id"`
//     User    *User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
//     DepartmentID uint `json:"department_id,omitempty"` // Your 'additional ID' for the department
//     // ...
//     Bio         string  `gorm:"type:text" json:"bio,omitempty"`
//     Image       string  `gorm:"size:255" json:"image,omitempty"`
//     Experience  int     `json:"experience,omitempty"`                  // years
//     Subjects    string  `gorm:"size:255" json:"subjects,omitempty"`    // comma-separated
//     Qualifications string `gorm:"size:255" json:"qualifications,omitempty"`

//     // Derived / transient fields (runtime calculation)
//     // TotalStudents int     `gorm:"-" json:"total_students,omitempty"`   // from enrollments
//     // Rating      float32 `gorm:"-" json:"rating,omitempty"`

//     // Relations
//     Courses     []Course     `gorm:"foreignKey:TeacherID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"courses,omitempty"`
//     Assignments []Assignment `gorm:"foreignKey:TeacherID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"assignments,omitempty"`

//     CreatedAt time.Time `json:"created_at"`
//     UpdatedAt time.Time `json:"updated_at"`
// }
type TeacherProfile struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID uint `gorm:"uniqueIndex;not null" json:"user_id"`

	// Links to User table (assuming User struct exists)
	User *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`

	// Departmental Authority Link (Nullable)
	DepartmentID *uint `json:"department_id,omitempty"` 
	Department *Department `gorm:"foreignKey:DepartmentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"department,omitempty"`

	// Teacher-specific fields (adjust these to match your existing fields)
	Bio  string `gorm:"type:text" json:"bio,omitempty"`
	Experience int `json:"experience,omitempty"`

	// NOTE: We REMOVE the Courses []Course reverse association here to prevent FK conflict.

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TeacherProfile) TableName() string {
	return "teacher_profiles"
}
