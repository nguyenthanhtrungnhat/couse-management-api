package models

import "github.com/google/uuid"

type Course struct {
	BaseModel

	InstructorID uuid.UUID `gorm:"type:uuid;index;not null"`
	Instructor   User

	CategoryID uuid.UUID `gorm:"type:uuid;index;not null"`
	Category   Category

	Title string `gorm:"size:255;not null"`

	Slug string `gorm:"size:255;uniqueIndex;not null"`

	Description *string

	ThumbnailURL *string

	PreviewVideoURL *string

	Price float64 `gorm:"type:numeric(12,2);default:0"`

	Status string `gorm:"size:20;default:draft"`

	AverageRating float64 `gorm:"type:numeric(3,2);default:0"`

	TotalStudents int `gorm:"default:0"`

	Sections []CourseSection `gorm:"foreignKey:CourseID"`

	Reviews []Review `gorm:"foreignKey:CourseID"`
}