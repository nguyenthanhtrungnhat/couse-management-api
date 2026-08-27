package seeders

import (
	"course-management/config"
	"course-management/models"
	"log"
)

func SeedReviews() {
	var student models.User
	var courses []models.Course

	if err := config.DB.
		Where("email = ?", "student@example.com").
		First(&student).Error; err != nil {
		return
	}

	if err := config.DB.Find(&courses).Error; err != nil {
		return
	}

	for _, course := range courses {

		var existing models.Review

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

		review := models.Review{
			UserID:   student.ID,
			CourseID: course.ID,
			Rating:   5,
			Comment:  "Very useful course. The explanations are clear and practical.",
		}

		if err := config.DB.Create(&review).Error; err != nil {
			log.Printf("❌ Review: %v", err)
			continue
		}

		log.Printf("✅ Review created")
	}
}