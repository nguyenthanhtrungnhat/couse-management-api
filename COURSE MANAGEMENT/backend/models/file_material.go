package models

import "github.com/google/uuid"

type FileMaterial struct {
	BaseModel

	LessonID uuid.UUID `json:"lesson_id"`
	Lesson   Lesson   `json:"lesson,omitempty"`

	FileName string `json:"file_name"`
	FileURL  string `json:"file_url"`
	FileType string `json:"file_type"`
	FileSize int64  `json:"file_size"`
}