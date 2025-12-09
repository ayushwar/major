package routes

import (
    "github.com/ayushwar/major/controllers"
    "github.com/ayushwar/major/middlewares"

    "github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all API endpoints
func RegisterRoutes(router *gin.Engine) {

    // User routes
    userRoutes := router.Group("/users")
    {
        userRoutes.POST("/register", controllers.Register)
        userRoutes.POST("/verify_email", controllers.VerifyEmail)
        userRoutes.POST("/login", controllers.Login)
        userRoutes.POST("/forget_password", controllers.ForgotPassword)
        userRoutes.POST("/reset_password", controllers.ResetPassword)
    }

    // Other resource routes
    CourseRoutes(router)
    // LectureRoutes(router)
    AssignmentRoutes(router)
    QuestionRoutes(router)
    OptionRoutes(router)
    RegisterEnrollmentRoutes(router)
    SubmissionRoutes(router)
    ProgressRoutes(router)
    CertificateRoutes(router)
    DepartmentRoutes(router)
}


func CourseRoutes(router *gin.Engine) {
    courses := router.Group("/courses")
    {
        // Public routes
        courses.GET("/", controllers.GetAllCourses)
        courses.GET("/:id", controllers.GetCourseByID)

        // Protected: teacher/admin only
        courses.POST("/",
            middlewares.AuthMiddleware(),
            middlewares.RoleMiddleware("teacher", "admin"),
            controllers.CreateCourse,
        )
        courses.PUT("/:id",
            middlewares.AuthMiddleware(),
            middlewares.RoleMiddleware("teacher", "admin"),
            controllers.UpdateCourse,
        )
        courses.DELETE("/:id",
            middlewares.AuthMiddleware(),
            middlewares.RoleMiddleware("teacher", "admin"),
            controllers.DeleteCourse,
        )
    }
}

// func LectureRoutes(router *gin.Engine) {
//     // Public: list lectures by course
//     router.GET("/courses/:id/lectures", controllers.GetLecturesByCourse)

//     // Protected: teacher/admin only
//     lectures := router.Group("/lectures")
//     lectures.Use(middlewares.AuthMiddleware(), middlewares.RoleMiddleware("teacher", "admin"))
//     {
//         lectures.POST("/", controllers.CreateLecture)
//         lectures.PUT("/:id", controllers.UpdateLecture)
//         lectures.DELETE("/:id", controllers.DeleteLecture)
//     }
// }

func AssignmentRoutes(router *gin.Engine) {
    assignments := router.Group("/assignments")
    {
        // Public: view assignments
        assignments.GET("/", controllers.GetAllAssignments)
        assignments.GET("/:id", controllers.GetAssignmentByID)

        // Protected: teacher/admin only for modification
        assignments.Use(middlewares.AuthMiddleware(), middlewares.RoleMiddleware("teacher", "admin"))
        {
            assignments.POST("/", controllers.CreateAssignment)
            assignments.PUT("/:id", controllers.UpdateAssignment)
            assignments.DELETE("/:id", controllers.DeleteAssignment)
        }
    }
}
func QuestionRoutes(router *gin.Engine) {
    // Assignment related questions — use :id consistently
    questions := router.Group("/assignments/:id/questions")
    {
        questions.GET("/", controllers.GetQuestionsByAssignment)

        questions.POST("/",
            middlewares.AuthMiddleware(),
            middlewares.RoleMiddleware("teacher", "admin"),
            controllers.CreateQuestion,
        )
    }

    // Individual question update/delete with :question_id param unchanged
    q := router.Group("/questions")
    q.Use(middlewares.AuthMiddleware(), middlewares.RoleMiddleware("teacher", "admin"))
    {
        q.PUT("/:question_id", controllers.UpdateQuestion)
        q.DELETE("/:question_id", controllers.DeleteQuestion)
    }
}


func OptionRoutes(router *gin.Engine) {
    // Question-related options, param :question_id with unique option id param :option_id
    options := router.Group("/questions/:question_id/options")
    {
        // Public: list options of question
        options.GET("/", controllers.GetOptions)

        // Protected: teacher/admin modify
        options.Use(middlewares.AuthMiddleware(), middlewares.RoleMiddleware("teacher", "admin"))
        {
            options.POST("/", controllers.CreateOption)
            options.PUT("/:option_id", controllers.UpdateOption)
            options.DELETE("/:option_id", controllers.DeleteOption)
        }
    }
}

func RegisterEnrollmentRoutes(r *gin.Engine) {
    // Enrollment routes protected for authenticated users
    enrollments := r.Group("/enrollments")
    enrollments.Use(middlewares.AuthMiddleware())

    // Student: enroll course
    enrollments.POST("", middlewares.RoleMiddleware("student"), controllers.EnrollCourse)

    // Student: get own enrollments
    r.GET("/users/:id/enrollments", middlewares.AuthMiddleware(), controllers.GetEnrollmentsByUser)

    // Teacher: get enrollments for course
    r.GET("/courses/:id/enrollments", middlewares.AuthMiddleware(), controllers.GetEnrollmentsByCourse)

    // Admin: update/delete enrollments
    enrollments.PUT("/:id", middlewares.RoleMiddleware("admin"), controllers.UpdateEnrollment)
    enrollments.DELETE("/:id", middlewares.RoleMiddleware("admin"), controllers.DeleteEnrollment)
}

func SubmissionRoutes(router *gin.Engine) {
    submissions := router.Group("/submissions")
    submissions.Use(middlewares.AuthMiddleware())
    {
        // Students: submit assignments
        submissions.POST("/", controllers.SubmitAssignment)

        // Students: fetch own submissions
        submissions.GET("/user/:id", controllers.GetSubmissionsByUser)

        // Teachers/Admin: fetch all submissions for assignment
        submissions.GET("/assignment/:id", controllers.GetSubmissionsByAssignment)
    }
}

func ProgressRoutes(router *gin.Engine) {
    progress := router.Group("/progress")
    progress.Use(middlewares.AuthMiddleware())
    {
        // Students: update progress
        progress.POST("/update", controllers.UpdateProgress)

        // Students/teachers/admin: fetch progress with params
        progress.GET("/:userId/:courseId", controllers.GetProgress)
    }
}

func CertificateRoutes(router *gin.Engine) {
    certs := router.Group("/certificates")
    certs.Use(middlewares.AuthMiddleware())
    {
        // Students: request certificate after completion
        certs.POST("/issue", controllers.IssueCertificate)

        // Students: get own certificates
        certs.GET("/user/:id", controllers.GetCertificatesByUser)

        // Admin/teacher: verify certificate by ID
        certs.GET("/:id", controllers.GetCertificateByID)
    }
}
func DepartmentRoutes(router *gin.Engine) {
	departments := router.Group("/departments")
	{
		// 1. Public Routes (Read Access for everyone)
		departments.GET("/", controllers.GetDepartments)
		departments.GET("/:id", controllers.GetDepartmentByID)

		// 2. Protected Routes (Admin Only for CUD operations)
		
		// Create Department (Admin Only)
		departments.POST("/",
			middlewares.AuthMiddleware(),
			middlewares.RoleMiddleware("admin"), // <-- Only 'admin' can create
			controllers.CreateDepartment,
		)

		// Update Department (Admin Only)
		departments.PUT("/:id",
			middlewares.AuthMiddleware(),
			middlewares.RoleMiddleware("admin"), // <-- Only 'admin' can update
			controllers.UpdateDepartment,
		)

		// Delete Department (Admin Only)
		departments.DELETE("/:id",
			middlewares.AuthMiddleware(),
			middlewares.RoleMiddleware("admin"), // <-- Only 'admin' can delete
			controllers.DeleteDepartment,
		)
	}
}