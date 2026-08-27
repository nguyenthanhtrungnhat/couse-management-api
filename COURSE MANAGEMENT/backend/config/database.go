package config

import (
	"fmt"
	"log"
	"os"

	"course-management/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect PostgreSQL: %v", err)
	}

	DB = db

	log.Println("✅ Connected to PostgreSQL!")

	// Auto Migrate
	err = DB.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Category{},
		&models.Course{},
		&models.CourseSection{},
		&models.Lesson{},
		&models.FileMaterial{},
		&models.Enrollment{},
		&models.LessonProgress{},
		&models.Review{},
		&models.Comment{},
		&models.Payment{},
		&models.PaymentLog{},
	)

	if err != nil {
		log.Fatalf("❌ AutoMigrate failed: %v", err)
	}

	log.Println("✅ Database migrated successfully!")
}
