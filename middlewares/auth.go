package middlewares

import (
	"errors"
	// "fmt"
	"os"
	// "strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4" // Check your JWT package version
)

// JWT_SECRET environment variable se secret key load karna
var secretKey = []byte(os.Getenv("JWT_SECRET"))



// GenerateToken creates JWT token with userID and role claims
func GenerateToken(userID uint, role string) (string, error) {
	// Note: We are using "userID" (string) for consistency and clearer retrieval
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID, // user_id (uint) will be stored
		"role":role,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString(secretKey)
}

// VerifyToken validates JWT token string and returns claims if valid
func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		// Token parsing failed (e.g., signature mismatch, expired)
		return nil, err
	}
    
    if !token.Valid {
        // Token is not valid (e.g., expiry time passed)
        return nil, errors.New("token is not valid")
    }

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, errors.New("invalid token claims format")
}

// AuthMiddleware protects routes with JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(401, gin.H{"error": "Authorization header missing"})
			ctx.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			ctx.JSON(401, gin.H{"error": "Invalid Authorization header format. Must be 'Bearer [token]'"})
			ctx.Abort()
			return
		}

		claims, err := VerifyToken(tokenParts[1])
		if err != nil {
			// Include details if available for debugging
			ctx.JSON(401, gin.H{"error": "Invalid or expired token", "details": err.Error()})
			ctx.Abort()
			return
		}

		// --- FIXES APPLIED HERE ---
		
		// 1. User ID (critical fix: type assertion from float64 to uint)
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
            // Fails if claim is missing or not a number
			ctx.JSON(401, gin.H{"error": "Token claim 'user_id' is missing or invalid format"})
			ctx.Abort()
			return
		}
        // Store as a uint for consistent use in application logic
		ctx.Set("userID", uint(userIDFloat))

		// 2. Role (Ensure it's a string)
		roleString, ok := claims["role"].(string)
		if !ok {
			ctx.JSON(401, gin.H{"error": "Token claim 'role' is missing or invalid format"})
			ctx.Abort()
			return
		}
		ctx.Set("role", roleString)

		ctx.Next()
	}
}

// RoleMiddleware checks if the user has required roles
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. Role retrieval
		roleValue, exists := ctx.Get("role")
		if !exists {
			// This should ideally never happen if AuthMiddleware ran successfully
			ctx.JSON(403, gin.H{"error": "Authorization error: user role not found in context"})
			ctx.Abort()
			return
		}

		// 2. Role type assertion (should be safe as AuthMiddleware sets it as string)
		userRole, ok := roleValue.(string)
		if !ok {
			ctx.JSON(403, gin.H{"error": "Internal error: invalid role type"})
			ctx.Abort()
			return
		}

		// 3. Check role against allowed roles
		for _, r := range allowedRoles {
			if userRole == r {
				ctx.Next() // Role matched, proceed
				return
			}
		}

		// 4. Forbidden if no role matches
		ctx.JSON(403, gin.H{"error": "Forbidden: insufficient permissions"})
		ctx.Abort()
	}
}