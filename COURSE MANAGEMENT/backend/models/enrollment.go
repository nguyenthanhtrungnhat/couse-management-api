package models

import "github.com/google/uuid"

type Enrollment struct {
	BaseModel

	UserID   uuid.UUID `json:"user_id"`
	CourseID uuid.UUID `json:"course_id"`

	User   User   `json:"user,omitempty"`
	Course Course `json:"course,omitempty"`

	Progress []LessonProgress `json:"progress,omitempty"`
}