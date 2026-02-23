package controllers

import (
	"github.com/ayushwar/major/video_repository"
	"github.com/gin-gonic/gin"
)


type ListHandler struct {
	Repo *videorepository.VideoRepository
}
func (h *ListHandler) List(c *gin.Context) {
	videos, err := h.Repo.List()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch videos"})
		return
	}

	result := make([]gin.H, 0, len(videos))
	for _, v := range videos {
		result = append(result, gin.H{
			"title":     v.Title,
			"embed_url": "https://www.youtube.com/embed/" + v.YouTubeVideoID,
		})
	}

	c.JSON(200, result)
}