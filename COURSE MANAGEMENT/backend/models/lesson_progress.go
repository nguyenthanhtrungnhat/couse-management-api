package models

import "github.com/google/uuid"

type LessonProgress struct {
	BaseModel

	EnrollmentID uuid.UUID `gorm:"type:uuid;index;not null"`

	LessonID uuid.UUID `gorm:"type:uuid;index;not null"`

	Enrollment Enrollment

	Lesson Lesson

	Completed bool

	WatchedSeconds int
}
