package controllers

import (
	"fmt"
	"time"

	"github.com/ayushwar/major/database"
	"github.com/ayushwar/major/models"
	"github.com/gin-gonic/gin"
)

// -----------------------------
// EnrollCourse → POST /enrollments
// -----------------------------
func EnrollCourse(ctx *gin.Context) {
	var input struct {
		CourseID uint `json:"course_id" binding:"required"`
	}

	// logged-in user from JWT
	userIDFromToken, _ := ctx.Get("userID")
	role, _ := ctx.Get("role")

	// only students can enroll themselves
	if role != "student" {
		ctx.JSON(403, gin.H{"error": "only students can enroll in courses"})
		return
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	// Check if course exists
	var course models.Course
	if err := database.DB.First(&course, input.CourseID).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "course not found"})
		return
	}

	// Check if already enrolled
	var existing models.Enrollment
	if err := database.DB.Where("user_id = ? AND course_id = ?", userIDFromToken, input.CourseID).
		First(&existing).Error; err == nil {
		ctx.JSON(400, gin.H{"error": "already enrolled in this course"})
		return
	}

	// Create enrollment
	enrollment := models.Enrollment{
		UserID:     userIDFromToken.(uint),
		CourseID:   input.CourseID,
		Progress:   0,
		Status:     "active",
		EnrolledAt: time.Now(),
	}

	if err := database.DB.Create(&enrollment).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to create enrollment", "details": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "enrolled successfully", "enrollment": enrollment})
}

// -----------------------------
// GetEnrollmentsByUser → GET /users/:id/enrollments
// -----------------------------
func GetEnrollmentsByUser(ctx *gin.Context) {
	userIDParam := ctx.Param("id")
	userIDFromToken, _ := ctx.Get("userID")
	role, _ := ctx.Get("role")

	// student → can only fetch their own enrollments
	if role == "student" && userIDParam != fmt.Sprint(userIDFromToken) {
		ctx.JSON(403, gin.H{"error": "students can only view their own enrollments"})
		return
	}

	var enrollments []models.Enrollment
	if err := database.DB.Preload("Course").
		Where("user_id = ?", userIDParam).
		Find(&enrollments).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to fetch enrollments", "details": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"enrollments": enrollments})
}

// -----------------------------
// GetEnrollmentsByCourse → GET /courses/:id/enrollments
// (teacher/admin)
// -----------------------------
func GetEnrollmentsByCourse(ctx *gin.Context) {
	courseID := ctx.Param("id")
	userIDFromToken, _ := ctx.Get("userID")
	role, _ := ctx.Get("role")

	// if teacher → ensure they own the course
	if role == "teacher" {
		var course models.Course
		if err := database.DB.First(&course, courseID).Error; err != nil {
			ctx.JSON(404, gin.H{"error": "course not found"})
			return
		}
		if course.TeacherID != userIDFromToken {
			ctx.JSON(403, gin.H{"error": "you can only view enrollments for your own courses"})
			return
		}
	}

	var enrollments []models.Enrollment
	if err := database.DB.Preload("User").
		Where("course_id = ?", courseID).
		Find(&enrollments).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to fetch enrollments", "details": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"enrollments": enrollments})
}

// -----------------------------
// UpdateEnrollment → PUT /enrollments/:id (admin only)
// -----------------------------
func UpdateEnrollment(ctx *gin.Context) {
	role, _ := ctx.Get("role")
	if role != "admin" {
		ctx.JSON(403, gin.H{"error": "only admins can update enrollments"})
		return
	}

	id := ctx.Param("id")
	var enrollment models.Enrollment

	if err := database.DB.First(&enrollment, id).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "enrollment not found"})
		return
	}

	var input struct {
		Progress      *float32   `json:"progress"`
		Status        *string    `json:"status"`
		CompletedAt   *time.Time `json:"completed_at"`
		CertificateID *string    `json:"certificate_id"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	if input.Progress != nil {
		enrollment.Progress = *input.Progress
	}
	if input.Status != nil {
		enrollment.Status = *input.Status
	}
	if input.CompletedAt != nil {
		enrollment.CompletedAt = input.CompletedAt
	}
	if input.CertificateID != nil {
		enrollment.CertificateID = input.CertificateID
	}

	if err := database.DB.Save(&enrollment).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to update enrollment", "details": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "enrollment updated successfully", "enrollment": enrollment})
}

// -----------------------------
// DeleteEnrollment → DELETE /enrollments/:id (admin only)
// -----------------------------
func DeleteEnrollment(ctx *gin.Context) {
	role, _ := ctx.Get("role")
	if role != "admin" {
		ctx.JSON(403, gin.H{"error": "only admins can delete enrollments"})
		return
	}

	id := ctx.Param("id")
	var enrollment models.Enrollment

	if err := database.DB.First(&enrollment, id).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "enrollment not found"})
		return
	}

	if err := database.DB.Delete(&enrollment).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to delete enrollment", "details": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "enrollment deleted successfully"})
}
