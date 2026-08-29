
package seeders

import (
	"context"
	"log"

	"course-management/config"

	"github.com/google/uuid"
)

func SeedFileMaterials() {
	ctx := context.Background()

	rows, err := config.DB.Query(
		ctx,
		`SELECT id FROM lessons ORDER BY created_at ASC`,
	)
	if err != nil {
		log.Printf("❌ Cannot load lessons: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var lessonID uuid.UUID

		if err := rows.Scan(&lessonID); err != nil {
			log.Printf("❌ Cannot read lesson ID: %v", err)
			continue
		}

		materials := []struct {
			fileName string
			fileURL  string
			fileType string
			fileSize int64
		}{
			{
				fileName: "lesson-notes.pdf",
				fileURL:  "https://example.com/files/lesson-notes.pdf",
				fileType: "application/pdf",
				fileSize: 102400,
			},
			{
				fileName: "source-code.zip",
				fileURL:  "https://example.com/files/source-code.zip",
				fileType: "application/zip",
				fileSize: 204800,
			},
		}

		for _, material := range materials {
			_, err := config.DB.Exec(
				ctx,
				`INSERT INTO file_materials (
					id,
					lesson_id,
					file_name,
					file_url,
					file_type,
					file_size,
					created_at,
					updated_at
				)
				VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`,
				uuid.New(),
				lessonID,
				material.fileName,
				material.fileURL,
				material.fileType,
				material.fileSize,
			)

			if err != nil {
				log.Printf(
					"❌ Material %s: %v",
					material.fileName,
					err,
				)
				continue
			}

			log.Printf("✅ Material: %s", material.fileName)
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("❌ Lesson iteration error: %v", err)
	}
}

