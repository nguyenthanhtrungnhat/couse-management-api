package seeders

import (
	"course-management/config"
	"course-management/models"
	"log"
)

func SeedCourseSections() {
	var courses []models.Course

	if err := config.DB.Find(&courses).Error; err != nil {
		log.Println("❌ Cannot load courses")
		return
	}

	for _, course := range courses {

		sections := []models.CourseSection{
			{
				CourseID:  course.ID,
				Title:     "Introduction",
				SortOrder: 1,
			},
			{
				CourseID:  course.ID,
				Title:     "Core Concepts",
				SortOrder: 2,
			},
			{
				CourseID:  course.ID,
				Title:     "Practical Project",
				SortOrder: 3,
			},
		}

		for _, section := range sections {
			var existing models.CourseSection

			err := config.DB.
				Where(
					"course_id = ? AND title = ?",
					section.CourseID,
					section.Title,
				).
				First(&existing).Error

			if err == nil {
				continue
			}

			if err := config.DB.Create(&section).Error; err != nil {
				log.Printf("❌ Section %s: %v", section.Title, err)
				continue
			}

			log.Printf("✅ Section: %s", section.Title)
		}
	}
}
