package models

import "github.com/google/uuid"

type CourseSection struct {
	BaseModel

	CourseID uuid.UUID `gorm:"type:uuid;index;not null"`
	Course   Course

	Title string `gorm:"size:255;not null"`

	SortOrder int

	Lessons []Lesson `gorm:"foreignKey:SectionID"`
}
