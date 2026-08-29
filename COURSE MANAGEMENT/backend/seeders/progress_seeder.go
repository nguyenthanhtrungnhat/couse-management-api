package seeders

import (
	"context"
	"log"

	"course-management/config"
	"course-management/models"

	"github.com/google/uuid"
)

// SeedProgress creates sample lesson progress records.
func SeedProgress() {
	log.Println("🌱 Seeding lesson progress...")

	ctx := context.Background()

	// Get all enrollments.
	rows, err := config.DB.Query(
		ctx,
		`
		SELECT id, course_id
		FROM enrollments
		ORDER BY created_at
		`,
	)
	if err != nil {
		log.Printf("❌ Failed to load enrollments: %v", err)
		return
	}
	defer rows.Close()

	var enrollments []models.Enrollment

	for rows.Next() {
		var enrollment models.Enrollment

		if err := rows.Scan(
			&enrollment.ID,
			&enrollment.CourseID,
		); err != nil {
			log.Printf("❌ Failed to scan enrollment: %v", err)
			return
		}

		enrollments = append(enrollments, enrollment)
	}

	if err := rows.Err(); err != nil {
		log.Printf("❌ Failed while reading enrollments: %v", err)
		return
	}

	if len(enrollments) == 0 {
		log.Println("⚠️ No enrollments found, skipping progress seeding")
		return
	}

	for _, enrollment := range enrollments {

		// Get all lessons belonging to the enrolled course.
		lessonRows, err := config.DB.Query(
			ctx,
			`
			SELECT
				l.id
			FROM lessons l
			INNER JOIN course_sections cs
				ON cs.id = l.section_id
			WHERE cs.course_id = $1
			ORDER BY cs.sort_order, l.sort_order
			`,
			enrollment.CourseID,
		)

		if err != nil {
			log.Printf(
				"❌ Failed to load lessons for course %s: %v",
				enrollment.CourseID,
				err,
			)
			continue
		}

		var lessonIDs []uuid.UUID

		for lessonRows.Next() {
			var lessonID uuid.UUID

			if err := lessonRows.Scan(&lessonID); err != nil {
				log.Printf(
					"❌ Failed to scan lesson: %v",
					err,
				)
				break
			}

			lessonIDs = append(lessonIDs, lessonID)
		}

		lessonRows.Close()

		if len(lessonIDs) == 0 {
			log.Printf(
				"⚠️ No lessons found for course %s",
				enrollment.CourseID,
			)
			continue
		}

		for index, lessonID := range lessonIDs {

			var watchedSeconds int64
			var completed bool

			switch index {
			case 0:
				// First lesson: completed.
				watchedSeconds = 600
				completed = true

			case 1:
				// Second lesson: partially watched.
				watchedSeconds = 300
				completed = false

			default:
				// Remaining lessons: not started.
				watchedSeconds = 0
				completed = false
			}

			progressID := uuid.New()

			_, err := config.DB.Exec(
				ctx,
				`
				INSERT INTO lesson_progresses (
					id,
					enrollment_id,
					lesson_id,
					watched_seconds,
					completed,
					created_at,
					updated_at
				)
				VALUES (
					$1,
					$2,
					$3,
					$4,
					$5,
					NOW(),
					NOW()
				)
				ON CONFLICT (enrollment_id, lesson_id)
				DO UPDATE SET
					watched_seconds = EXCLUDED.watched_seconds,
					completed = EXCLUDED.completed,
					updated_at = NOW()
				`,
				progressID,
				enrollment.ID,
				lessonID,
				watchedSeconds,
				completed,
			)

			if err != nil {
				log.Printf(
					"❌ Progress for lesson %s: %v",
					lessonID,
					err,
				)
				continue
			}
		}

		log.Printf(
			"✅ Progress seeded for enrollment %s",
			enrollment.ID,
		)
	}

	log.Println("✅ Lesson progress seeding completed")
}
