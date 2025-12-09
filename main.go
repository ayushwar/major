package main

import (
	"log"

	"github.com/ayushwar/major/database"
	"github.com/ayushwar/major/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Add this init function to your main.go file
func init() {
	// Ise sabse pehle load karein taki JWT_SECRET available ho
	if err := godotenv.Load(); err != nil {
		// Yeh line sirf warning degi agar .env file nahi milti.
		log.Println("NOTE: No .env file found or unable to load.")
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		// Yeh line sirf warning degi agar .env file nahi milti.
		log.Println("NOTE: No .env file found or unable to load.")
	}

	database.ConnectDB()

	server := gin.Default()

	routes.RegisterRoutes(server)

	port := ":8080"
	log.Println("Server running on http://localhost" + port)
	if err := server.Run(port); err != nil {
		log.Fatal(" Failed to start server: ", err)
	}
}
