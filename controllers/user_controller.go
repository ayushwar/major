package controllers

import (
	"fmt"

	"time"

	"github.com/ayushwar/major/database"
	"github.com/ayushwar/major/middlewares"
	"github.com/ayushwar/major/models"
	"github.com/ayushwar/major/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// PendingUser holds temporary registration data (not persisted until verified)
type PendingUser struct {
	User         models.User
	OTP          string
	OTPExpiresAt time.Time
	Password     string
}

// in-memory temporary store for unverified users
// Note: this is ephemeral and will be lost on server restart. Use Redis for production.
var pendingUsers = make(map[string]PendingUser)

// Register: collect user info, generate OTP and store in-memory (NOT in DB yet)
func Register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}
	if !utils.ValidateEmail(user.Email) {
		ctx.JSON(400, gin.H{"error": "invalid email format"})
		return
	}

	// Check if email already exists in DB (only verified users are in DB)
	var existingUser models.User
	if err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		ctx.JSON(400, gin.H{"error": "email already exists"})
		return
	}
	plainTextPassword := user.Password

	// If there's already a pending registration for this email, overwrite with a new OTP
	otp, err := utils.GenerateOTP()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "failed to generate otp"})
		return
	}

	pendingUsers[user.Email] = PendingUser{
		User:         user,
		OTP:          otp,
		OTPExpiresAt: time.Now().Add(5 * time.Minute),
		Password:     plainTextPassword,
	}

	// Send OTP email
	subject := "Verify Your Email - OTP"
	body := fmt.Sprintf("Hello %s,\n\nYour OTP for email verification is: %s\nThis OTP will expire in 5 minutes.\n\nThanks!", user.Name, otp)
	if err := utils.SendEmail(user.Email, subject, body); err != nil {
		// If email fails, remove pending entry to avoid stale records
		delete(pendingUsers, user.Email)
		ctx.JSON(500, gin.H{"error": "failed to send OTP email", "details": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "OTP sent, please verify email"})
}

// VerifyEmail: validate OTP from pendingUsers, then persist user to DB
func VerifyEmail(ctx *gin.Context) {
	type VerifyInput struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}
	var input VerifyInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	pending, ok := pendingUsers[input.Email]
	if !ok {
		ctx.JSON(400, gin.H{"error": "user not found or not registered yet"})
		return
	}

	// validate OTP and expiry
	if pending.OTP != input.OTP || time.Now().After(pending.OTPExpiresAt) {
		ctx.JSON(400, gin.H{"error": "invalid or expired OTP"})
		return
	}
	fmt.Println("Pending plain password before hashing:", pending.Password)

	// --- HASHING and FINAL DB PREP ---

	// 1. Plaintext Password ko check karein
	plainPassword := pending.Password
	if plainPassword == "" {
		fmt.Println("ERROR: Plaintext password is empty for:", input.Email)
		ctx.JSON(500, gin.H{"error": "registration error: password data missing"})
		return
	}
	fmt.Println("DEBUG: Pending plain password before hashing (length):", len(plainPassword))

	// 2. Password ko hash karein
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("ERROR: Failed to hash password:", err)
		ctx.JSON(500, gin.H{"error": "failed to hash password", "details": err.Error()})
		return
	}
	// Hashed password ko check karein (length 60 honi chahiye)
	fmt.Println("DEBUG: After hashing (Hash string length):", len(hashed))
	// **NOTE: Is line mein aapko hash string terminal mein dikhai degi**
	fmt.Println("DEBUG: Final Hashed Password to be stored:", string(hashed))

	// 3. User Model ko update karein (Local Copy ka use karke)
	// Hum pending.User ko ek naye variable mein copy kar rahe hain
	verifiedUser := pending.User
	verifiedUser.Password = string(hashed) // Hashed password set karein
	verifiedUser.IsVerified = true         // Verified flag set karein

	// 4. Database mein save karein
	if err := database.DB.Create(&verifiedUser).Error; err != nil {
		fmt.Println("ERROR: Failed to save user to DB:", err)
		ctx.JSON(500, gin.H{"error": "failed to save user", "details": err.Error()})
		return
	}
	fmt.Println("SUCCESS: User saved to DB:", verifiedUser.Email)

	// 5. Cleanup pending entry
	delete(pendingUsers, input.Email)

	ctx.JSON(200, gin.H{"message": "email verified successfully, user registered"})
}

// Login handles user login and JWT issuing
func Login(ctx *gin.Context) {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var input LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		ctx.JSON(400, gin.H{"error": "invalid email or password"})
		return
	}

	if !user.IsVerified {
		ctx.JSON(401, gin.H{"error": "email is not verified"})
		return
	}
	// fmt.Println("DB stored password:", user.Password)
	// fmt.Println("User entered password:", input.Password)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		ctx.JSON(401, gin.H{"error": "invalid password"})
		fmt.Println("DB stored password:", user.Password)
		fmt.Println("User entered password:", input.Password)
		return
	}

	tokenString, err := middlewares.GenerateToken(user.ID, user.Role)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "failed to generate token", "details": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "login successfully", "token": tokenString})
}

// ForgotPassword handles password reset token generation and email sending
// ForgotPassword sends OTP to user email instead of reset token
func ForgotPassword(ctx *gin.Context) {
	type Input struct {
		Email string `json:"email" binding:"required,email"`
	}

	var input Input
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		ctx.JSON(400, gin.H{"error": "user not found"})
		return
	}

	// generate OTP
	otp, err := utils.GenerateOTP()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "failed to generate otp"})
		return
	}

	user.ResetToken = otp                              // reuse ResetToken column for OTP
	user.ResetExpiry = time.Now().Add(5 * time.Minute) // OTP valid for 5 minutes
	if err := database.DB.Save(&user).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to save otp", "details": err.Error()})
		return
	}

	subject := "Password Reset OTP"
	body := fmt.Sprintf(
		"Hello %s,\n\nYour OTP to reset your password is: %s\nThis OTP will expire in 5 minutes.\n\nThanks!",
		user.Name, otp,
	)

	if err := utils.SendEmail(user.Email, subject, body); err != nil {
		ctx.JSON(500, gin.H{"error": "failed to send OTP email", "details": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "OTP sent to email"})
}

// ResetPassword verifies OTP and updates the password
func ResetPassword(ctx *gin.Context) {
	type Input struct {
		Email       string `json:"email" binding:"required,email"`
		OTP         string `json:"otp" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	var input Input
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400, gin.H{"error": "invalid request, all fields required", "details": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		ctx.JSON(400, gin.H{"error": "invalid email"})
		return
	}

	// validate OTP and expiry
	if user.ResetToken != input.OTP || time.Now().After(user.ResetExpiry) {
		ctx.JSON(400, gin.H{"error": "invalid or expired OTP"})
		return
	}

	// hash new password
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "failed to hash password", "details": err.Error()})
		return
	}

	user.Password = string(hashed)
	user.ResetToken = ""           // clear OTP
	user.ResetExpiry = time.Time{} // clear expiry

	if err := database.DB.Save(&user).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to update password", "details": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "password reset successfully"})
}
