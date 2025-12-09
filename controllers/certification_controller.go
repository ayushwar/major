package controllers

import (
	"time"

	"github.com/ayushwar/major/database"
	"github.com/ayushwar/major/models"
	"github.com/ayushwar/major/utils"
	"github.com/gin-gonic/gin"
)

// IssueCertificate → POST /certificates/issue
// Auto-issue only when course progress is 100%
func IssueCertificate(ctx *gin.Context) {
	var req struct {
		UserID   uint `json:"user_id"`
		CourseID uint `json:"course_id"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	// Check enrollment
	var enrollment models.Enrollment
	if err := database.DB.Where("user_id = ? AND course_id = ?", req.UserID, req.CourseID).
		First(&enrollment).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "enrollment not found"})
		return
	}

	// Ensure course is fully completed
	if enrollment.Progress < 100 {
		ctx.JSON(400, gin.H{"error": "course not completed, certificate cannot be issued"})
		return
	}

	// Create certificate
	cert := models.Certificate{
		UserID:   req.UserID,
		CourseID: req.CourseID,
		IssuedAt: time.Now(),
		CertCode: utils.GenerateCertificateCode(), // util function
	}

	if err := database.DB.Create(&cert).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to issue certificate"})
		return
	}

	ctx.JSON(201, gin.H{"message": "certificate issued successfully", "certificate": cert})
}

// GetCertificatesByUser → GET /certificates/user/:id
func GetCertificatesByUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	var certs []models.Certificate

	if err := database.DB.Where("user_id = ?", userID).Find(&certs).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to fetch certificates"})
		return
	}

	ctx.JSON(200, gin.H{"certificates": certs})
}

// GetCertificateByID → GET /certificates/:id
func GetCertificateByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var cert models.Certificate

	if err := database.DB.First(&cert, id).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "certificate not found"})
		return
	}

	ctx.JSON(200, gin.H{"certificate": cert})
}
