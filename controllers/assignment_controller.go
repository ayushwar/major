package controllers

import (
	
	"strconv"

	"github.com/ayushwar/major/database"
	"github.com/ayushwar/major/models"
	"github.com/gin-gonic/gin"
)// CreateAssignment → POST /assignments
func CreateAssignment(ctx *gin.Context) {
	var assignment models.Assignment
	if err := ctx.ShouldBindJSON(&assignment); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&assignment).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create assignment", "details": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Assignment created successfully", "assignment": assignment})
}

// GetAllAssignments → GET /assignments
func GetAllAssignments(ctx *gin.Context) {
	var assignments []models.Assignment
	if err := database.DB.Preload("Questions.Options").Find(&assignments).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to fetch assignments", "details": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"assignments": assignments})
}

// GetAssignmentByID → GET /assignments/:id
func GetAssignmentByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var assignment models.Assignment
	if err := database.DB.Preload("Questions.Options").First(&assignment, id).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "Assignment not found"})
		return
	}

	ctx.JSON(200, gin.H{"assignment": assignment})
}

// UpdateAssignment → PUT /assignments/:id
func UpdateAssignment(ctx *gin.Context) {
	id := ctx.Param("id")
	var assignment models.Assignment

	if err := database.DB.First(&assignment, id).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "Assignment not found"})
		return
	}

	var input models.Assignment
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	assignment.Title = input.Title
	assignment.Description = input.Description
	assignment.CourseID = input.CourseID

	if err := database.DB.Save(&assignment).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to update assignment", "details": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Assignment updated successfully", "assignment": assignment})
}

// DeleteAssignment → DELETE /assignments/:id
func DeleteAssignment(ctx *gin.Context) {
	id := ctx.Param("id")
	assignmentID, _ := strconv.Atoi(id)

	if err := database.DB.Delete(&models.Assignment{}, assignmentID).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to delete assignment", "details": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Assignment deleted successfully"})
}