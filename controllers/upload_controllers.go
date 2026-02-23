package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ayushwar/major/models"
	"github.com/ayushwar/major/services"
	videorepository "github.com/ayushwar/major/video_repository"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	Repo *videorepository.VideoRepository
}

func parseUint(s string) uint {
	v, _ := strconv.ParseUint(s, 10, 64)
	return uint(v)
}

func (h *UploadHandler) Upload(c *gin.Context) {
	// ---------- Form values ----------
	title := c.PostForm("title")
	desc := c.PostForm("description")
	batchID := c.PostForm("batch_id")

	if title == "" || batchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title and batch_id required"})
		return
	}

	// ---------- Auth ----------
	userID := c.GetUint("user_id") // set by AuthMiddleware

	// ---------- File ----------
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "video file required"})
		return
	}

	// ---------- Save temporarily ----------
	_ = os.MkdirAll("uploads", 0755)
	path := filepath.Join("uploads", file.Filename)

	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	defer os.Remove(path)

	// ---------- YouTube upload ----------
	videoID, err := services.UploadToYouTube(
		path,
		title,
		desc,
		os.Getenv("YOUTUBE_CLIENT_ID"),
		os.Getenv("YOUTUBE_CLIENT_SECRET"),
		os.Getenv("YOUTUBE_REFRESH_TOKEN"),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ---------- Save DB ----------
	err = h.Repo.Create(&models.Video{
		Title:          title,
		Description:    desc,
		YouTubeVideoID: videoID,
		BatchID:        parseUint(batchID),
		TeacherID:      userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db insert failed"})
		return
	}

	// ---------- Response ----------
	c.JSON(http.StatusOK, gin.H{
		"message":   "Video uploaded successfully",
		"video_id":  videoID,
		"embed_url": "https://www.youtube.com/embed/" + videoID,
	})
}