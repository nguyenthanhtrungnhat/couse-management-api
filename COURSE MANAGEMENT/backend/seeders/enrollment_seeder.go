package seeders

import (
	"context"
	"log"

	"course-management/config"
	"course-management/models"

	"github.com/google/uuid"
)

func SeedEnrollments() {
	log.Println("🌱 Seeding enrollments...")

	// Get student
	var studentID uuid.UUID
	err := config.DB.QueryRow(
		context.Background(),
		`
		SELECT id
		FROM users
		WHERE email = $1
		LIMIT 1
		`,
		"student@example.com",
	).Scan(&studentID)

	if err != nil {
		log.Printf("❌ Enrollment: failed to find student: %v", err)
		return
	}

	// Get courses
	rows, err := config.DB.Query(
		context.Background(),
		`
		SELECT id
		FROM courses
		ORDER BY created_at
		`,
	)

	if err != nil {
		log.Printf("❌ Enrollment: failed to get courses: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var courseID uuid.UUID

		if err := rows.Scan(&courseID); err != nil {
			log.Printf("❌ Enrollment: failed to scan course: %v", err)
			continue
		}

		enrollmentID := uuid.New()

		_, err := config.DB.Exec(
			context.Background(),
			`
			INSERT INTO enrollments (
				id,
				created_at,
				updated_at,
				user_id,
				course_id
			)
			VALUES (
				$1,
				NOW(),
				NOW(),
				$2,
				$3
			)
			ON CONFLICT (user_id, course_id)
			DO NOTHING
			`,
			enrollmentID,
			studentID,
			courseID,
		)

		if err != nil {
			log.Printf(
				"❌ Enrollment: student %s -> course %s: %v",
				studentID,
				courseID,
				err,
			)
			continue
		}

		log.Printf(
			"✅ Enrollment: student %s -> course %s",
			studentID,
			courseID,
		)
	}

	if err := rows.Err(); err != nil {
		log.Printf("❌ Enrollment rows error: %v", err)
	}

	log.Println("✅ Enrollment seeding completed")
}

// Keep models import available if other seeders share this pattern.
var _ = models.Enrollment{}
