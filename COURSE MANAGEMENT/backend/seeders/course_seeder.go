package seeders

import (
	"context"
	"log"

	"course-management/config"

	"github.com/google/uuid"
)

func SeedCourses() {
	ctx := context.Background()

	var instructorID uuid.UUID
	var categoryID uuid.UUID

	err := config.DB.QueryRow(
		ctx,
		`SELECT id FROM users WHERE email = $1`,
		"instructor@example.com",
	).Scan(&instructorID)

	if err != nil {
		log.Printf("❌ Demo instructor not found: %v", err)
		return
	}

	err = config.DB.QueryRow(
		ctx,
		`SELECT id FROM categories WHERE name = $1`,
		"Programming",
	).Scan(&categoryID)

	if err != nil {
		log.Printf("❌ Programming category not found: %v", err)
		return
	}

	courses := []struct {
		title           string
		slug            string
		description     string
		thumbnailURL    string
		previewVideoURL string
		price           float64
		status          string
		averageRating   float64
		totalStudents   int
	}{
		{
			title:           "Go Fiber REST API",
			slug:            "go-fiber-rest-api",
			description:     "Build REST APIs using Go Fiber and PostgreSQL.",
			thumbnailURL:    "https://example.com/go-fiber.jpg",
			previewVideoURL: "https://example.com/go-fiber-preview.mp4",
			price:           199000,
			status:          "published",
			averageRating:   4.5,
			totalStudents:   2,
		},
		{
			title:           "PostgreSQL Database Fundamentals",
			slug:            "postgresql-database-fundamentals",
			description:     "Learn PostgreSQL, SQL queries and relational database design.",
			thumbnailURL:    "https://example.com/postgresql.jpg",
			previewVideoURL: "https://example.com/postgresql-preview.mp4",
			price:           149000,
			status:          "published",
			averageRating:   4.3,
			totalStudents:   1,
		},
	}

	for _, course := range courses {
		id := uuid.New()

		_, err := config.DB.Exec(
			ctx,
			`INSERT INTO courses (
				id,
				instructor_id,
				category_id,
				title,
				slug,
				description,
				thumbnail_url,
				preview_video_url,
				price,
				status,
				average_rating,
				total_students,
				created_at,
				updated_at
			)
			VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8,
				$9, $10, $11, $12, NOW(), NOW()
			)`,
			id,
			instructorID,
			categoryID,
			course.title,
			course.slug,
			course.description,
			course.thumbnailURL,
			course.previewVideoURL,
			course.price,
			course.status,
			course.averageRating,
			course.totalStudents,
		)

		if err != nil {
			log.Printf("❌ Course %s: %v", course.title, err)
			continue
		}

		log.Printf("✅ Course: %s", course.title)
	}
}

