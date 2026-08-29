
package seeders

import (
	"context"
	"log"

	"course-management/config"

	"github.com/google/uuid"
)

func SeedLessons() {
	ctx := context.Background()

	rows, err := config.DB.Query(
		ctx,
		`SELECT id FROM course_sections ORDER BY created_at ASC`,
	)
	if err != nil {
		log.Printf("❌ Cannot load sections: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var sectionID uuid.UUID

		if err := rows.Scan(&sectionID); err != nil {
			log.Printf("❌ Cannot read section ID: %v", err)
			continue
		}

		lessons := []struct {
			title           string
			videoURL        string
			durationSeconds int
			isPreview       bool
			sortOrder       int
		}{
			{
				title:           "Getting Started",
				videoURL:        "https://example.com/videos/getting-started.mp4",
				durationSeconds: 600,
				isPreview:       true,
				sortOrder:       1,
			},
			{
				title:           "Understanding the Basics",
				videoURL:        "https://example.com/videos/basics.mp4",
				durationSeconds: 900,
				isPreview:       false,
				sortOrder:       2,
			},
			{
				title:           "Building a Practical Example",
				videoURL:        "https://example.com/videos/practical.mp4",
				durationSeconds: 1200,
				isPreview:       false,
				sortOrder:       3,
			},
		}

		for _, lesson := range lessons {
			_, err := config.DB.Exec(
				ctx,
				`INSERT INTO lessons (
					id,
					section_id,
					title,
					video_url,
					duration_seconds,
					is_preview,
					sort_order,
					created_at,
					updated_at
				)
				VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())`,
				uuid.New(),
				sectionID,
				lesson.title,
				lesson.videoURL,
				lesson.durationSeconds,
				lesson.isPreview,
				lesson.sortOrder,
			)

			if err != nil {
				log.Printf(
					"❌ Lesson %s: %v",
					lesson.title,
					err,
				)
				continue
			}

			log.Printf("✅ Lesson: %s", lesson.title)
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("❌ Section iteration error: %v", err)
	}
}

