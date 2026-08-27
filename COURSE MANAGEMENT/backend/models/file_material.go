package models

import (
	"github.com/google/uuid"
)

type FileMaterial struct {
	BaseModel

	LessonID uuid.UUID `gorm:"type:uuid;not null;index"`
	Lesson   Lesson    `gorm:"foreignKey:LessonID"`

	FileName string `gorm:"size:255;not null"`
	FileURL  string `gorm:"type:text;not null"`
	FileType string `gorm:"size:100"`
	FileSize int64  `gorm:"not null;default:0"`
}
