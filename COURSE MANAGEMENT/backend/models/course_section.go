package models

import "github.com/google/uuid"

type CourseSection struct {
	BaseModel

	CourseID uuid.UUID `json:"course_id"`
	Course   Course   `json:"course,omitempty"`

	Title     string `json:"title"`
	SortOrder int    `json:"sort_order"`

	Lessons []Lesson `json:"lessons,omitempty"`
}