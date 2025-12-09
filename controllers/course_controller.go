package controllers

import (
	"fmt"
	"github.com/ayushwar/major/database"
	"github.com/ayushwar/major/models"
	"github.com/gin-gonic/gin"
	
)

// helper: get user role from context
func getUserRole(ctx *gin.Context) string {
	role, exists := ctx.Get("role")
	if !exists {
		return ""
	}
	// Note: Role is set as string in AuthMiddleware
	return role.(string)
}

func getContextUserID(ctx *gin.Context) (uint, bool) {
	userID, exists := ctx.Get("userID") // Check key 'userID' or 'user_id'
	if !exists {
		return 0, false
	}
	
	// Handle uint (expected) or float64 (common from JWT claims parsing)
	if id, ok := userID.(uint); ok {
		return id, true
	}
	if floatID, ok := userID.(float64); ok {
		return uint(floatID), true
	}
	return 0, false
}


// CreateCourse handles course creation with departmental authorization
func CreateCourse(ctx *gin.Context) {
	var course models.Course

	// 1. JSON Data Binding with validation
	if err := ctx.ShouldBindJSON(&course); err != nil {
		ctx.JSON(400, gin.H{
			"error":"invalid request data",
			"details": err.Error(),
		})
		return
	}

	// 2. Assign TeacherID from JWT (this is the User ID)
	teacherID, ok := getContextUserID(ctx)
	if !ok {
		ctx.JSON(401, gin.H{
			"error": "user not authenticated",
		})
		return
	}
	course.TeacherID = teacherID

	// 3. Validate DepartmentID is provided (mandatory check)
	if course.DepartmentID == 0 {
		ctx.JSON(400, gin.H{
			"error": "department_id is required",
		})
		return
	}

	// 4. DEPARTMENTAL AUTHORIZATION CHECK (Your Business Logic)
	var teacherProfile models.TeacherProfile
	
	// 4A. Fetch the Teacher's Profile using the UserID
	if err := database.DB.Where("user_id = ?", teacherID).First(&teacherProfile).Error; err != nil {
		ctx.JSON(403, gin.H{
			"error": "Teacher profile validation failed",
			"details": "You must have a complete teacher profile to create courses.",
		})
		return
	}

	// 4B. Check if teacher has an assigned DepartmentID
	if teacherProfile.DepartmentID == nil || *teacherProfile.DepartmentID == 0 {
		ctx.JSON(403, gin.H{
			"error": "Department not assigned",
			"details": "Your teacher profile must be assigned to a department.",
		})
		return
	}

	// 4C. Verify teacher's department matches the course's requested department
	if *teacherProfile.DepartmentID != course.DepartmentID {
		ctx.JSON(403, gin.H{
			"error": "Authorization Failed: Department Mismatch",
			"details": fmt.Sprintf(
				"You are authorized only for Department ID %d, not the requested Department ID %d",
				*teacherProfile.DepartmentID,
				course.DepartmentID,
			),
		})
		return
	}

	// 5. Verify the requested Department exists (Prevents FK error for Department)
	var department models.Department
	if err := database.DB.First(&department, course.DepartmentID).Error; err != nil {
		ctx.JSON(401, gin.H{
			"error": "invalid department_id",
			"details": "The specified department does not exist in the database.",
		})
		return
	}

	// 6. Create the course in the database
	if err := database.DB.Create(&course).Error; err != nil {
		// This should only show genuine errors now, not Error 1452 (if migrations are run)
		fmt.Println("DB ERROR during Course Creation:", err.Error())
		ctx.JSON(500, gin.H{
			"error" :"failed to create course",
			"details": err.Error(),
		})
		return
	}

	// 7. Load relationships for response (Optional, but good practice)
	// You might need to adjust the preloads based on your User and Department models
	database.DB.Preload("Teacher").Preload("Department").First(&course, course.ID)

	// Success response
	ctx.JSON(201, gin.H{
		"message": "course created successfully",
		"course": course,
	})
}

// GetAllCourses → GET /courses
func GetAllCourses(ctx *gin.Context) {
	var courses []models.Course
    // Preload Department and TeacherProfile data for complete view
	if err := database.DB.Preload("Department").Preload("Teacher").Find(&courses).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to fetch courses", "details": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"courses": courses})
}

// GetCourseByID → GET /courses/:id
func GetCourseByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var course models.Course
    // Preload data
	if err := database.DB.Preload("Department").Preload("Teacher").First(&course, id).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "course not found"})
		return
	}
	ctx.JSON(200, gin.H{"course": course})
}

// controllers/course_controller.go

// UpdateCourse → PUT /courses/:id
func UpdateCourse(ctx *gin.Context) {
    id := ctx.Param("id")
    var existingCourse models.Course

    // 1. Course Fetch karna (Fetch the existing course record)
    if err := database.DB.First(&existingCourse, id).Error; err != nil {
        ctx.JSON(404, gin.H{"error": "course not found"})
        return
    }

    // 2. JWT se User ID aur Role nikalna
    loggedUserID, ok := getContextUserID(ctx)
    role := getUserRole(ctx) // Assuming this helper exists and works

    if !ok || role == "" {
        ctx.JSON(401, gin.H{"error": "Authorization data missing"})
        ctx.Abort()
        return
    }

    // 2A. Teacher Role Check: Sirf Teacher aur Admin hi update kar sakte hain
    if role != "admin" && role != "teacher" {
        ctx.JSON(403, gin.H{"error": "Forbidden: Only authorized teachers or admins can update courses."})
        ctx.Abort()
        return
    }

    // 3. Ownership Check (CRITICAL)
    // Sirf Admin ya Course ka original Teacher hi update kar sakta hai.
    if role != "admin" && existingCourse.TeacherID != loggedUserID {
        ctx.JSON(403, gin.H{"error": "Forbidden: You can only update courses you created."})
        ctx.Abort()
        return
    }

    // 4. Input Bind karna
    var input models.Course
    if err := ctx.ShouldBindJSON(&input); err != nil {
        ctx.JSON(400, gin.H{"error": "invalid request", "details": err.Error()})
        return
    }
    
    // 5. DEPARTMENTAL AUTHORIZATION CHECK (Agar DepartmentID change ho raha hai)
    
    // Check if DepartmentID is provided in the input AND if it's different from the existing one.
    if input.DepartmentID != 0 && input.DepartmentID != existingCourse.DepartmentID {
        
        // 5A. Teacher ka Profile Fetch karna
        var teacherProfile models.TeacherProfile
        if err := database.DB.Where("user_id = ?", loggedUserID).First(&teacherProfile).Error; err != nil {
             ctx.JSON(403, gin.H{"error": "Profile error: Teacher profile not found for authorization check."})
             return
        }

        // 5B. Check if teacher is assigned any department
        if teacherProfile.DepartmentID == nil || *teacherProfile.DepartmentID == 0 {
            ctx.JSON(403, gin.H{"error": "Forbidden: Your profile must be assigned to a department to change course department."})
            return
        }

        // 5C. Verify teacher's assigned department matches the NEW requested department (input.DepartmentID)
        if *teacherProfile.DepartmentID != input.DepartmentID {
            ctx.JSON(403, gin.H{
                "error": "Authorization Failed: Department Mismatch",
                "details": fmt.Sprintf(
                    "You are authorized only for Department ID %d, but tried to change to Department ID %d.",
                    *teacherProfile.DepartmentID,
                    input.DepartmentID,
                ),
            })
            return
        }
        
        // Agar check pass ho gaya, tab DepartmentID ko update karein
        existingCourse.DepartmentID = input.DepartmentID
    }
    
    // 6. Fields Update karna (Use provided input fields)
    
    // Note: We use the existing field-by-field update logic for safety.
    if input.Title != "" {
        existingCourse.Title = input.Title
    }
    if input.Description != "" {
        existingCourse.Description = input.Description
    }
    if input.Code != "" {
        existingCourse.Code = input.Code
    }
    if input.Credits > 0 {
        existingCourse.Credits = input.Credits
    }
    // Update other fields as needed (e.g., Price, Duration, Level)
    
    // 7. Database Save karna
    if err := database.DB.Save(&existingCourse).Error; err != nil {
        ctx.JSON(500, gin.H{"error": "failed to update course", "details": err.Error()})
        return
    }

    // Load relationships for complete response
    database.DB.Preload("Teacher").Preload("Department").First(&existingCourse, existingCourse.ID)

    ctx.JSON(200, gin.H{"message": "course updated successfully", "course": existingCourse})
}

// DeleteCourse → DELETE /courses/:id
func DeleteCourse(ctx *gin.Context) {
	id := ctx.Param("id")
	var course models.Course

	// 1. Course Fetch karna
	if err := database.DB.First(&course, id).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "course not found or already deleted"})
		return
	}

	// 2. JWT se User ID aur Role nikalna
	loggedUserID, ok := getContextUserID(ctx)
	role := getUserRole(ctx)

	if !ok || role == "" {
		ctx.JSON(401, gin.H{"error": "Authorization data missing"})
		ctx.Abort()
		return
	}
	
	// 3. Ownership Check (CRITICAL)
	// Sirf Admin ya Course ka Teacher hi delete kar sakta hai.
	if role != "admin" && course.TeacherID != loggedUserID {
		ctx.JSON(403, gin.H{"error": "Forbidden: You can only delete courses you created."})
		ctx.Abort()
		return
	}
	
	// 4. Database Delete karna
	if err := database.DB.Delete(&course).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to delete course", "details": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "course deleted successfully"})
}