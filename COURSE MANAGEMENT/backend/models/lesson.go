
package models

import "github.com/google/uuid"

type Lesson struct {
	BaseModel

	SectionID uuid.UUID `json:"section_id"`

	Title           string  `json:"title"`
	VideoURL        *string `json:"video_url,omitempty"`
	DurationSeconds int64   `json:"duration_seconds"`
	IsPreview       bool    `json:"is_preview"`
	SortOrder       int64   `json:"sort_order"`

	Section       *CourseSection `json:"section,omitempty"`
	FileMaterials []FileMaterial  `json:"file_materials,omitempty"`
}

// TableName returns the database table name.
func (Lesson) TableName() string {
	return "lessons"
}

// NewLesson creates a new Lesson model.
func NewLesson(
	sectionID uuid.UUID,
	title string,
	videoURL *string,
	durationSeconds int64,
	isPreview bool,
	sortOrder int64,
) *Lesson {
	return &Lesson{
		BaseModel: BaseModel{
			ID: uuid.New(),
		},

		SectionID:       sectionID,
		Title:           title,
		VideoURL:        videoURL,
		DurationSeconds: durationSeconds,
		IsPreview:       isPreview,
		SortOrder:       sortOrder,
	}
}

