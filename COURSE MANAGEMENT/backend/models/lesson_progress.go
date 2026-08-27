package models

import "github.com/google/uuid"

type LessonProgress struct {
	BaseModel

	EnrollmentID uuid.UUID `gorm:"type:uuid;not null;index"`
	Enrollment   Enrollment `gorm:"foreignKey:EnrollmentID"`

	LessonID uuid.UUID `gorm:"type:uuid;not null;index"`
	Lesson   Lesson `gorm:"foreignKey:LessonID"`

	Completed      bool `gorm:"not null;default:false"`
	WatchedSeconds int  `gorm:"not null;default:0"`
}