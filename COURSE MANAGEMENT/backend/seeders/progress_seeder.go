package seeders

import (
	"course-management/config"
	"course-management/models"
	"log"
)

func SeedProgress() {
	var enrollments []models.Enrollment

	if err := config.DB.Find(&enrollments).Error; err != nil {
		log.Println("❌ Cannot load enrollments")
		return
	}

	for _, enrollment := range enrollments {

		var course models.Course

		if err := config.DB.
			Preload("Sections.Lessons").
			First(&course, "id = ?", enrollment.CourseID).Error; err != nil {
			continue
		}

		for _, section := range course.Sections {
			for _, lesson := range section.Lessons {

				var existing models.LessonProgress

				err := config.DB.
					Where(
						"enrollment_id = ? AND lesson_id = ?",
						enrollment.ID,
						lesson.ID,
					).
					First(&existing).Error

				if err == nil {
					continue
				}

				progress := models.LessonProgress{
					EnrollmentID:   enrollment.ID,
					LessonID:       lesson.ID,
					Completed:      lesson.SortOrder == 1,
					WatchedSeconds: 300,
				}

				if err := config.DB.Create(&progress).Error; err != nil {
					log.Printf("❌ Progress: %v", err)
				}
			}
		}
	}
}
