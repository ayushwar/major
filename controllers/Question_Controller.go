package controllers

import (
	"github.com/ayushwar/major/database"
	"github.com/ayushwar/major/models"
	"github.com/gin-gonic/gin"
)

// CreateQuestion → POST /assignments/:id/questions
func CreateQuestion(ctx *gin.Context)  {
	assignmentID:=ctx.Param("id")
	var input struct{
		Text string `json:"text" binding:"required"`
	}
	if err:=ctx.ShouldBindBodyWithJSON(&input);err!=nil{
		ctx.JSON(400,gin.H{"error":"invaild request"})
		return
	}
	
	// Ensure assignment exists
	var assignment models.Assignment
	if err := database.DB.First(&assignment, assignmentID).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "assignment not found"})
		return
	}

	question := models.Question{
		AssignmentID: assignment.ID,
		Text:         input.Text,
	}

	if err := database.DB.Create(&question).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to create question", "details": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "question created successfully", "question": question})


}
// GetQuestionsByAssignment → GET /assignments/:id/questions
func GetQuestionsByAssignment(ctx *gin.Context) {
	assignmentID := ctx.Param("id")

	var questions []models.Question
	if err := database.DB.Where("assignment_id = ?", assignmentID).Find(&questions).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to fetch questions", "details": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"questions": questions})
}

// UpdateQuestion → PUT /questions/:id
func UpdateQuestion(ctx *gin.Context) {
	id := ctx.Param("id")

	var question models.Question
	if err := database.DB.First(&question, id).Error; err != nil {
		ctx.JSON(400, gin.H{"error": "question not found"})
		return
	}

	var input struct {
		Text string `json:"text"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(404, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	if input.Text != "" {
		question.Text = input.Text
	}

	if err := database.DB.Save(&question).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to update question", "details": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "question updated successfully", "question": question})
}

// DeleteQuestion → DELETE /questions/:id
func DeleteQuestion(ctx *gin.Context) {
	id := ctx.Param("id")

	var question models.Question
	if err := database.DB.First(&question, id).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "question not found"})
		return
	}

	if err := database.DB.Delete(&question).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to delete question", "details": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "question deleted successfully"})
}

// ------------------------------------------------------------------
//							option controllers 
// ------------------------------------------------------------------
// CreateOption → add option to a question
func CreateOption(c *gin.Context) {
	var option models.Option
	if err := c.ShouldBindJSON(&option); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&option).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, option)
}

// GetOptions → get all options for a question
func GetOptions(c *gin.Context) {
	questionID := c.Param("question_id")
	var options []models.Option

	if err := database.DB.Where("question_id = ?", questionID).Find(&options).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, options)
}

// UpdateOption → edit option (e.g. text or correctness)
func UpdateOption(c *gin.Context) {
	id := c.Param("id")
	var option models.Option

	if err := database.DB.First(&option, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Option not found"})
		return
	}

	if err := c.ShouldBindJSON(&option); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&option).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, option)
}

// DeleteOption → delete an option
func DeleteOption(c *gin.Context) {
	id := c.Param("id")
	var option models.Option

	if err := database.DB.First(&option, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Option not found"})
		return
	}

	if err := database.DB.Delete(&option).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Option deleted successfully"})
}