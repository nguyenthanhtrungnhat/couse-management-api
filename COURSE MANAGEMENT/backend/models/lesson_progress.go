package models

import "github.com/google/uuid"

type LessonProgress struct {
	BaseModel

	EnrollmentID uuid.UUID  `json:"enrollment_id"`
	Enrollment   Enrollment `json:"enrollment,omitempty"`

	LessonID uuid.UUID `json:"lesson_id"`
	Lesson   Lesson   `json:"lesson,omitempty"`

	Completed      bool `json:"completed"`
	WatchedSeconds int  `json:"watched_seconds"`
}