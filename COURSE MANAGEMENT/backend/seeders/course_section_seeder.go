package seeders

import (
	"context"
	"log"

	"course-management/config"

	"github.com/google/uuid"
)

func SeedCourseSections() {
	ctx := context.Background()

	rows, err := config.DB.Query(
		ctx,
		`SELECT id FROM courses ORDER BY created_at ASC`,
	)
	if err != nil {
		log.Printf("❌ Cannot load courses: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var courseID uuid.UUID

		if err := rows.Scan(&courseID); err != nil {
			log.Printf("❌ Cannot read course ID: %v", err)
			continue
		}

		sections := []struct {
			title     string
			sortOrder int
		}{
			{
				title:     "Introduction",
				sortOrder: 1,
			},
			{
				title:     "Core Concepts",
				sortOrder: 2,
			},
			{
				title:     "Practical Project",
				sortOrder: 3,
			},
		}

		for _, section := range sections {
			_, err := config.DB.Exec(
				ctx,
				`INSERT INTO course_sections (
					id,
					course_id,
					title,
					sort_order,
					created_at,
					updated_at
				)
				VALUES ($1, $2, $3, $4, NOW(), NOW())`,
				uuid.New(),
				courseID,
				section.title,
				section.sortOrder,
			)

			if err != nil {
				log.Printf(
					"❌ Section %s: %v",
					section.title,
					err,
				)
				continue
			}

			log.Printf("✅ Section: %s", section.title)
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("❌ Course iteration error: %v", err)
	}
}
