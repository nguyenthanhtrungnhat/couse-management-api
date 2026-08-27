package seeders

import (
	"course-management/config"
	"course-management/models"
	"log"
)

func SeedEnrollments() {
	var student models.User
	var courses []models.Course

	if err := config.DB.
		Where("email = ?", "student@example.com").
		First(&student).Error; err != nil {
		log.Println("❌ Demo student not found")
		return
	}

	if err := config.DB.Find(&courses).Error; err != nil {
		return
	}

	for _, course := range courses {

		var existing models.Enrollment

		err := config.DB.
			Where(
				"user_id = ? AND course_id = ?",
				student.ID,
				course.ID,
			).
			First(&existing).Error

		if err == nil {
			continue
		}

		enrollment := models.Enrollment{
			UserID:   student.ID,
			CourseID: course.ID,
		}

		if err := config.DB.Create(&enrollment).Error; err != nil {
			log.Printf("❌ Enrollment: %v", err)
			continue
		}

		log.Printf("✅ Enrollment created")
	}
}
