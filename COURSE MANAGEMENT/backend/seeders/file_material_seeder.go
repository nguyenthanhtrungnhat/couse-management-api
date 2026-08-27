package seeders

import (
	"errors"
	"log"

	"course-management/config"
	"course-management/models"

	"gorm.io/gorm"
)

func SeedFileMaterials() {
	var lessons []models.Lesson

	if err := config.DB.Find(&lessons).Error; err != nil {
		log.Println("❌ Cannot load lessons:", err)
		return
	}

	for _, lesson := range lessons {

		materials := []models.FileMaterial{
			{
				LessonID: lesson.ID,
				FileName: "lesson-notes.pdf",
				FileURL:  "https://example.com/files/lesson-notes.pdf",
				FileType: "application/pdf",
				FileSize: 102400,
			},
			{
				LessonID: lesson.ID,
				FileName: "source-code.zip",
				FileURL:  "https://example.com/files/source-code.zip",
				FileType: "application/zip",
				FileSize: 204800,
			},
		}

		for _, material := range materials {

			var existing models.FileMaterial

			err := config.DB.
				Where(
					"lesson_id = ? AND file_name = ?",
					material.LessonID,
					material.FileName,
				).
				First(&existing).Error

			// Material already exists
			if err == nil {
				log.Printf(
					"⏭️ Material already exists: %s",
					material.FileName,
				)
				continue
			}

			// Unexpected database error
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf(
					"❌ Check material %s: %v",
					material.FileName,
					err,
				)
				continue
			}

			// Material does not exist → create
			if err := config.DB.Create(&material).Error; err != nil {
				log.Printf(
					"❌ Material %s: %v",
					material.FileName,
					err,
				)
				continue
			}

			log.Printf(
				"✅ Material created: %s",
				material.FileName,
			)
		}
	}
}
