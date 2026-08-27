package seeders

import (
	"course-management/config"
	"course-management/models"
	"log"
)

func SeedFileMaterials() {
	var lessons []models.Lesson

	if err := config.DB.Find(&lessons).Error; err != nil {
		log.Println("❌ Cannot load lessons")
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

			if err == nil {
				continue
			}

			if err := config.DB.Create(&material).Error; err != nil {
				log.Printf("❌ Material %s: %v", material.FileName, err)
				continue
			}

			log.Printf("✅ Material: %s", material.FileName)
		}
	}
}
