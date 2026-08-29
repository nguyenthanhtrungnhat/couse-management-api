package seeders

import (
	"context"
	"log"

	"course-management/config"

	"github.com/google/uuid"
)

func SeedProgress() {
	ctx := context.Background()

	rows, err := config.DB.Query(
		ctx,
		`SELECT id, user_id, course_id
		 FROM enrollments
		 ORDER BY created_at ASC`,
	)
	if err != nil {
		log.Printf("❌ Cannot load enrollments: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			enrollmentID uuid.UUID
			userID       uuid.UUID
			courseID     uuid.UUID
		)

		if err := rows.Scan(
			&enrollmentID,
			&userID,
			&courseID,
		); err != nil {
			log.Printf("❌ Cannot read enrollment: %v", err)
			continue
		}

		lessonRows, err := config.DB.Query(
			ctx,
			`SELECT l.id
			 FROM lessons l
			 JOIN course_sections cs
			   ON cs.id = l.section_id
			 WHERE cs.course_id = $1
			 ORDER BY cs.sort_order, l.sort_order`,
			courseID,
		)

		if err != nil {
			log.Printf(
				"❌ Cannot load lessons for course %s: %v",
				courseID,
				err,
			)
			continue
		}

		var lessonIDs []uuid.UUID

		for lessonRows.Next() {
			var lessonID uuid.UUID

			if err := lessonRows.Scan(&lessonID); err != nil {
				log.Printf("❌ Cannot read lesson ID: %v", err)
				continue
			}

			lessonIDs = append(lessonIDs, lessonID)
		}

		lessonRows.Close()

		if err := lessonRows.Err(); err != nil {
			log.Printf("❌ Lesson iteration error: %v", err)
			continue
		}

		for i, lessonID := range lessonIDs {
			completed := i == 0

			_, err := config.DB.Exec(
				ctx,
				`INSERT INTO lesson_progress (
					id,
					user_id,
					lesson_id,
					completed,
					created_at,
					updated_at
				)
				VALUES ($1, $2, $3, $4, NOW(), NOW())`,
				uuid.New(),
				userID,
				lessonID,
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

			log.Printf(
				"✅ Progress: user %s -> lesson %s",
				userID,
				lessonID,
			)
		}

		log.Printf(
			"✅ Progress seeded for enrollment %s",
			enrollmentID,
		)
	}

	if err := rows.Err(); err != nil {
		log.Printf("❌ Enrollment iteration error: %v", err)
	}
}
