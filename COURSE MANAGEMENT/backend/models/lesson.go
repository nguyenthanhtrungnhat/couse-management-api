package models

import "github.com/google/uuid"

type Lesson struct {
	BaseModel

	SectionID uuid.UUID     `json:"section_id"`
	Section   CourseSection `json:"section,omitempty"`

	Title           string `json:"title"`
	VideoURL        string `json:"video_url"`
	DurationSeconds int    `json:"duration_seconds"`
	IsPreview       bool   `json:"is_preview"`
	SortOrder       int    `json:"sort_order"`

	FileMaterials []FileMaterial `json:"file_materials,omitempty"`
}