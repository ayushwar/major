package controllers

import (
	"github.com/ayushwar/major/database"
	

	"github.com/ayushwar/major/models"
	"github.com/gin-gonic/gin"
	
)


// Create Department
func CreateDepartment(c *gin.Context) {
	var dept models.Department

	if err := c.ShouldBindJSON(&dept); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&dept).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, dept)
}

// Get All Departments
func GetDepartments(c *gin.Context) {
	var departments []models.Department

	if err := database.DB.Find(&departments).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, departments)
}

// Get Department by ID (with Courses)
func GetDepartmentByID(c *gin.Context) {
	id := c.Param("id")

	var dept models.Department

	if err := database.DB.Preload("Courses").First(&dept, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "department not found"})
		return
	}

	c.JSON(200, dept)
}

// Update Department
func UpdateDepartment(c *gin.Context) {
	id := c.Param("id")

	var dept models.Department

	if err := database.DB.First(&dept, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "department not found"})
		return
	}

	var updateData models.Department
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&dept).Updates(updateData)

	c.JSON(200, dept)
}

// Delete Department
func DeleteDepartment(c *gin.Context) {
	id := c.Param("id")

	var dept models.Department

	if err := database.DB.First(&dept, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "department not found"})
		return
	}

	database.DB.Delete(&dept)

	c.JSON(200, gin.H{"message": "department deleted"})
}
