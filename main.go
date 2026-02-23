package main

import (
	"log"
	"os"

	"github.com/ayushwar/major/controllers"
	"github.com/ayushwar/major/database"
	"github.com/ayushwar/major/routes"
	videorepository "github.com/ayushwar/major/video_repository"
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

	// ---- VIDEO MODULE INITIALIZATION (MERGED HERE) ----
	videoRepo := videorepository.NewVideoRepository(database.DB)

	uploadHandler := &controllers.UploadHandler{
		Repo: videoRepo,
	}
	// uploadHandler.Config.ClientID = os.Getenv("YT_CLIENT_ID")
	// uploadHandler.Config.ClientSecret = os.Getenv("YT_CLIENT_SECRET")
	// uploadHandler.Config.RefreshToken = os.Getenv("YT_REFRESH_TOKEN")

	listHandler := &controllers.ListHandler{
		Repo: videoRepo,
	}
	// --------------------------------------------------

	// Register all routes (main + merged video routes)
	routes.RegisterRoutes(server, uploadHandler, listHandler)

	port := os.Getenv("PORT")
		if port == "" {
		port = "8080" // local fallback only
	}
	log.Println("Server running on http://localhost" + port)

	if err := server.Run(":"+port); err != nil {
		log.Fatal(" Failed to start server: ", err)
	}
}
