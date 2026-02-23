package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	DatabaseURL         string
	YouTubeClientID     string
	YouTubeClientSecret string
	YouTubeRefreshToken string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:                os.Getenv("PORT"),
		// DatabaseURL:         os.Getenv("DATABASE_URL"),
		YouTubeClientID:     os.Getenv("YOUTUBE_CLIENT_ID"),
		YouTubeClientSecret: os.Getenv("YOUTUBE_CLIENT_SECRET"),
		YouTubeRefreshToken: os.Getenv("YOUTUBE_REFRESH_TOKEN"),
	}

	if cfg.YouTubeClientID == "" || cfg.YouTubeClientSecret == "" || cfg.YouTubeRefreshToken == "" {
		log.Fatal("YouTube OAuth credentials are missing in .env")
	}

	return cfg
}