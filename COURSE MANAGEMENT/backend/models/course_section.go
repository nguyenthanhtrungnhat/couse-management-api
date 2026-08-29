
package models

import (
	"github.com/google/uuid"
)

type CourseSection struct {
	BaseModel

	CourseID uuid.UUID `json:"course_id"`

	Title     string `json:"title"`
	SortOrder int64  `json:"sort_order"`

	Course   *Course `json:"course,omitempty"`
	Lessons  []Lesson `json:"lessons,omitempty"`
}

// TableName returns the database table name.
func (CourseSection) TableName() string {
	return "course_sections"
}

// NewCourseSection creates a new course section.
func NewCourseSection(
	courseID uuid.UUID,
	title string,
	sortOrder int64,
) *CourseSection {
	return &CourseSection{
		BaseModel: BaseModel{
			ID: uuid.New(),
		},

		CourseID:  courseID,
		Title:     title,
		SortOrder: sortOrder,
	}
}

