package seeders

import (
	"course-management/config"
	"course-management/models"
	"log"
)

func SeedComments() {
	var student models.User
	var instructor models.User
	var course models.Course

	if err := config.DB.
		Where("email = ?", "student@example.com").
		First(&student).Error; err != nil {
		return
	}

	if err := config.DB.
		Where("email = ?", "instructor@example.com").
		First(&instructor).Error; err != nil {
		return
	}

	if err := config.DB.
		Order("created_at ASC").
		First(&course).Error; err != nil {
		return
	}

	var existing models.Comment

	err := config.DB.
		Where(
			"user_id = ? AND course_id = ? AND content = ?",
			student.ID,
			course.ID,
			"Can you explain the authentication part in more detail?",
		).
		First(&existing).Error

	if err != nil {
		comment := models.Comment{
			UserID:   student.ID,
			CourseID: course.ID,
			Content:  "Can you explain the authentication part in more detail?",
		}

		if err := config.DB.Create(&comment).Error; err != nil {
			log.Printf("❌ Comment: %v", err)
		}
	}

	// Instructor reply
	var parent models.Comment

	if err := config.DB.
		Where(
			"user_id = ? AND course_id = ?",
			student.ID,
			course.ID,
		).
		First(&parent).Error; err != nil {
		return
	}

	var reply models.Comment

	if err := config.DB.
		Where(
			"user_id = ? AND parent_id = ?",
			instructor.ID,
			parent.ID,
		).
		First(&reply).Error; err == nil {
		return
	}

	reply = models.Comment{
		UserID:   instructor.ID,
		CourseID: course.ID,
		ParentID: &parent.ID,
		Content:  "Sure. We will cover JWT authentication in the next lesson.",
	}

	if err := config.DB.Create(&reply).Error; err != nil {
		log.Printf("❌ Reply: %v", err)
	}
}
