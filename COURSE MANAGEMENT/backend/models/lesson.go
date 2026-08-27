package models

import (
	"github.com/google/uuid"
)

type Lesson struct {
	BaseModel

	SectionID uuid.UUID     `gorm:"type:uuid;not null;index"`
	Section   CourseSection `gorm:"foreignKey:SectionID"`

	Title           string `gorm:"size:255;not null"`
	VideoURL        string `gorm:"type:text;not null"`
	DurationSeconds int    `gorm:"not null;default:0"`
	IsPreview       bool   `gorm:"not null;default:false"`
	SortOrder       int    `gorm:"not null;default:0"`

	FileMaterials []FileMaterial `gorm:"foreignKey:LessonID"`
}
