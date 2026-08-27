package lesson

import (
	"time"

	"course-management/dto/filematerial"
)

type LessonResponse struct {
	ID              string                               `json:"id"`
	SectionID       string                               `json:"section_id"`
	Title           string                               `json:"title"`
	VideoURL        string                               `json:"video_url"`
	DurationSeconds int                                  `json:"duration_seconds"`
	IsPreview       bool                                 `json:"is_preview"`
	SortOrder       int                                  `json:"sort_order"`
	CreatedAt       time.Time                            `json:"created_at"`
	UpdatedAt       time.Time                            `json:"updated_at"`

	FileMaterials []filematerial.FileMaterialResponse `json:"file_materials"`
}