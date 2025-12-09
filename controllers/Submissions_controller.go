package controllers

import (
	
	"time"

	"github.com/ayushwar/major/database"
	"github.com/ayushwar/major/models"
	"github.com/gin-gonic/gin"
)

// SubmitAssignment → POST /submissions
func SubmitAssignment(ctx *gin.Context) {
	var req struct {
		AssignmentID uint            `json:"assignment_id"`
		UserID       uint            `json:"user_id"`
		Answers      map[uint]uint   `json:"answers"` // QuestionID → SelectedOptionID
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	// Fetch assignment questions with options
	var questions []models.Question
	if err := database.DB.Preload("Options").Where("assignment_id = ?", req.AssignmentID).Find(&questions).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to fetch questions"})
		return
	}

	// Calculate score
	score := 0
	for _, q := range questions {
		correctOption := uint(0)
		for _, opt := range q.Options {
			if opt.IsCorrect {
				correctOption = opt.ID
				break
			}
		}

		if req.Answers[q.ID] == correctOption {
			score++
		}
	}

	// Save submission
	submission := models.Submission{
		AssignmentID: req.AssignmentID,
		UserID:       req.UserID,
		Score:        score,
		SubmittedAt:  time.Now(),
	}

	if err := database.DB.Create(&submission).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to save submission"})
		return
	}

	ctx.JSON(200, gin.H{
		"message":    "submission saved successfully",
		"score":      score,
		"submission": submission,
	})
}

// GetSubmissionsByUser → GET /submissions/user/:id
func GetSubmissionsByUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	var submissions []models.Submission

	if err := database.DB.Where("user_id = ?", userID).Find(&submissions).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to fetch submissions"})
		return
	}

	ctx.JSON(200, gin.H{"submissions": submissions})
}

// GetSubmissionsByAssignment → GET /submissions/assignment/:id
func GetSubmissionsByAssignment(ctx *gin.Context) {
	assignmentID := ctx.Param("id")
	var submissions []models.Submission

	if err := database.DB.Where("assignment_id = ?", assignmentID).Find(&submissions).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to fetch submissions"})
		return
	}

	ctx.JSON(200, gin.H{"submissions": submissions})
}
