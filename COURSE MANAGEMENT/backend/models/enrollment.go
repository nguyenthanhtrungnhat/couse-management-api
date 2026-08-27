package models

import "github.com/google/uuid"

type Enrollment struct {
	BaseModel

	UserID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_user_course;not null"`

	CourseID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_user_course;not null"`

	User User

	Course Course

	Progress []LessonProgress
}
