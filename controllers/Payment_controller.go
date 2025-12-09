package controllers

import (
	"time"

	"github.com/ayushwar/major/database"
	"github.com/ayushwar/major/models"
	"github.com/gin-gonic/gin"
)

// CreatePayment - User initiates payment
func CreatePayment(ctx *gin.Context) {
	var payment models.Payment

	if err := ctx.ShouldBindJSON(&payment); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Step 1: Fetch profile of the user
	var profile models.Profile
	if err := database.DB.Where("user_id = ?", payment.UserID).First(&profile).Error; err == nil {
		// Step 2: Check if student is verified
		if profile.Verified {
			// Example: 20% discount
			discount := payment.Amount * 0.20
			payment.DiscountApplied = discount
			payment.Amount = payment.Amount - discount
		}
	}

	// Step 3: Set default values
	payment.Status = "pending"
	payment.CreatedAt = time.Now()

	// Step 4: Save payment
	if err := database.DB.Create(&payment).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to create payment"})
		return
	}

	ctx.JSON(201, gin.H{
		"message": "payment created successfully",
		"payment": payment,
	})
}



