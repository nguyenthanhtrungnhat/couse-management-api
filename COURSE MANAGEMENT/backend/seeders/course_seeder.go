package seeders

import (
	"course-management/config"
	"course-management/models"
	"log"
)

func SeedCourses() {
	var instructor models.User
	var category models.Category

	if err := config.DB.
		Where("email = ?", "instructor@example.com").
		First(&instructor).Error; err != nil {
		log.Println("❌ Demo instructor not found")
		return
	}

	if err := config.DB.
		Where("name = ?", "Programming").
		First(&category).Error; err != nil {
		log.Println("❌ Programming category not found")
		return
	}

	courses := []models.Course{
		{
			InstructorID:    instructor.ID,
			CategoryID:      category.ID,
			Title:           "Go Fiber REST API",
			Slug:            "go-fiber-rest-api",
			Description:     stringPtr("Build REST APIs using Go Fiber and PostgreSQL."),
			ThumbnailURL:    stringPtr("https://example.com/go-fiber.jpg"),
			PreviewVideoURL: stringPtr("https://example.com/go-fiber-preview.mp4"),
			Price:           199000,
			Status:          "published",
			AverageRating:   4.5,
			TotalStudents:   2,
		},
		{
			InstructorID:    instructor.ID,
			CategoryID:      category.ID,
			Title:           "PostgreSQL Database Fundamentals",
			Slug:            "postgresql-database-fundamentals",
			Description:     stringPtr("Learn PostgreSQL, SQL queries and relational database design."),
			ThumbnailURL:    stringPtr("https://example.com/postgresql.jpg"),
			PreviewVideoURL: stringPtr("https://example.com/postgresql-preview.mp4"),
			Price:           149000,
			Status:          "published",
			AverageRating:   4.3,
			TotalStudents:   1,
		},
	}

	for _, course := range courses {
		var existing models.Course

		if err := config.DB.
			Where("slug = ?", course.Slug).
			First(&existing).Error; err == nil {
			continue
		}

		if err := config.DB.Create(&course).Error; err != nil {
			log.Printf("❌ Course %s: %v", course.Title, err)
			continue
		}

		log.Printf("✅ Course: %s", course.Title)
	}
}
