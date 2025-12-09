package controllers

import (
	"github.com/ayushwar/major/database"
	"github.com/ayushwar/major/models"
	"github.com/gin-gonic/gin"
)

// UpdateProgress → POST /progress/update
func UpdateProgress(ctx *gin.Context) {
	var req struct {
		UserID   uint `json:"user_id"`
		CourseID uint `json:"course_id"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	// Count total assignments in the course
	var totalAssignments int64
	if err := database.DB.Model(&models.Assignment{}).
		Where("course_id = ?", req.CourseID).
		Count(&totalAssignments).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to count assignments"})
		return
	}

	// Count how many assignments this user submitted in this course
	var completedAssignments int64
	if err := database.DB.Model(&models.Submission{}).
		Joins("JOIN assignments ON submissions.assignment_id = assignments.id").
		Where("submissions.user_id = ? AND assignments.course_id = ?", req.UserID, req.CourseID).
		Count(&completedAssignments).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to count submissions"})
		return
	}

	// Calculate percentage progress
	progress := float32(0)
	if totalAssignments > 0 {
		progress = (float32(completedAssignments) / float32(totalAssignments)) * 100
	}

	// Update enrollment record
	var enrollment models.Enrollment
	if err := database.DB.Where("user_id = ? AND course_id = ?", req.UserID, req.CourseID).
		First(&enrollment).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "enrollment not found"})
		return
	}
	enrollment.Progress = progress

	if err := database.DB.Save(&enrollment).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to update progress"})
		return
	}

	ctx.JSON(200, gin.H{
		"message":       "progress updated successfully",
		"progress":      progress,
		"total":         totalAssignments,
		"completed":     completedAssignments,
		"enrollment_id": enrollment.ID,
	})
}

// GetProgress → GET /progress/:userId/:courseId
func GetProgress(ctx *gin.Context) {
	userID := ctx.Param("userId")
	courseID := ctx.Param("courseId")

	var enrollment models.Enrollment
	if err := database.DB.Where("user_id = ? AND course_id = ?", userID, courseID).
		First(&enrollment).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "enrollment not found"})
		return
	}

	ctx.JSON(200, gin.H{
		"user_id":   userID,
		"course_id": courseID,
		"progress":  enrollment.Progress,
	})
}
