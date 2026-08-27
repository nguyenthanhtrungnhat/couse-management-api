package seeders

import (
	"course-management/config"
	"course-management/models"
	"log"
)

func SeedLessons() {
	var sections []models.CourseSection

	if err := config.DB.Find(&sections).Error; err != nil {
		log.Println("❌ Cannot load sections")
		return
	}

	for _, section := range sections {

		lessons := []models.Lesson{
			{
				SectionID:       section.ID,
				Title:           "Getting Started",
				VideoURL:        "https://example.com/videos/getting-started.mp4",
				DurationSeconds: 600,
				IsPreview:       true,
				SortOrder:       1,
			},
			{
				SectionID:       section.ID,
				Title:           "Understanding the Basics",
				VideoURL:        "https://example.com/videos/basics.mp4",
				DurationSeconds: 900,
				IsPreview:       false,
				SortOrder:       2,
			},
			{
				SectionID:       section.ID,
				Title:           "Building a Practical Example",
				VideoURL:        "https://example.com/videos/practical.mp4",
				DurationSeconds: 1200,
				IsPreview:       false,
				SortOrder:       3,
			},
		}

		for _, lesson := range lessons {

			var existing models.Lesson

			err := config.DB.
				Where(
					"section_id = ? AND title = ?",
					lesson.SectionID,
					lesson.Title,
				).
				First(&existing).Error

			if err == nil {
				continue
			}

			if err := config.DB.Create(&lesson).Error; err != nil {
				log.Printf("❌ Lesson %s: %v", lesson.Title, err)
				continue
			}

			log.Printf("✅ Lesson: %s", lesson.Title)
		}
	}
}
