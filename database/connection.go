package database

import (
	"fmt"
	"log"
	"os"

	"github.com/ayushwar/major/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB initializes the database connection
func ConnectDB() {
	// Load environment variables from .env (optional if already set in system)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using system environment variables")
	}

	// Read env vars
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	// Validate (fail fast if any are missing)
	if user == "" || pass == "" || host == "" || port == "" || name == "" {
		log.Fatal("❌ Missing required database environment variables")
	}

	// Build DSN
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name,
	)

	// Connect
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database: ", err)
	}

	DB = db
	log.Println("✅ Database connected!")

	// Run migrations
	err = db.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Course{},
		&models.Lecture{},
		&models.Enrollment{},
		&models.Assignment{},
		&models.Submission{},
		&models.Payment{},
		&models.CollegeVerification{},
		&models.Progress{},
		&models.Certificate{},
		&models.Department{},
	)
	if err != nil {
		log.Fatal("❌ Migration failed: ", err)
	}

	log.Println("✅ All models migrated successfully!")
}
