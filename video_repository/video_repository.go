package videorepository

import (
	"github.com/ayushwar/major/models"
	"gorm.io/gorm"
)

type VideoRepository struct {    
    DB *gorm.DB
}
// ✅ CONSTRUCTOR (THIS IS WHAT YOU WERE MISSING)
func NewVideoRepository(db *gorm.DB) *VideoRepository {
	return &VideoRepository{DB: db}
}
// Create using GORM
func (r *VideoRepository) Create(video *models.Video) error {
    // GORM handles the INSERT query automatically
    return r.DB.Create(video).Error
}

// List using GORM
func (r *VideoRepository) List() ([]models.Video, error) {
    var videos []models.Video
    // GORM handles the SELECT * and scanning into the slice
    err := r.DB.Order("created_at desc").Find(&videos).Error
    return videos, err
}