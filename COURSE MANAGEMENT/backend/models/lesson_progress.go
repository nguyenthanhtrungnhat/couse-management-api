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

func (LessonProgress) TableName() string {
	return "lesson_progresses"
}

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
