
package models

import "github.com/google/uuid"

type LessonProgress struct {
	BaseModel

	EnrollmentID uuid.UUID `json:"enrollment_id"`
	LessonID     uuid.UUID `json:"lesson_id"`

	Completed      bool  `json:"completed"`
	WatchedSeconds int64 `json:"watched_seconds"`

	Enrollment *Enrollment `json:"enrollment,omitempty"`
	Lesson     *Lesson     `json:"lesson,omitempty"`
}

// TableName returns the database table name.
func (LessonProgress) TableName() string {
	return "lesson_progresses"
}

// NewLessonProgress creates a new lesson progress record.
func NewLessonProgress(
	enrollmentID uuid.UUID,
	lessonID uuid.UUID,
) *LessonProgress {
	return &LessonProgress{
		BaseModel: BaseModel{
			ID: uuid.New(),
		},

		EnrollmentID:   enrollmentID,
		LessonID:       lessonID,
		Completed:      false,
		WatchedSeconds: 0,
	}
}

