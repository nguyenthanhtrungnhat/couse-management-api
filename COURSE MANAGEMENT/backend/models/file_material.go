
package models

import "github.com/google/uuid"

type FileMaterial struct {
	BaseModel

	LessonID uuid.UUID `json:"lesson_id"`

	FileName string  `json:"file_name"`
	FileURL  *string `json:"file_url,omitempty"`
	FileSize int64   `json:"file_size"`
	FileType *string `json:"file_type,omitempty"`
	SortOrder int64  `json:"sort_order"`

	Lesson *Lesson `json:"lesson,omitempty"`
}

// TableName returns the database table name.
func (FileMaterial) TableName() string {
	return "file_materials"
}

// NewFileMaterial creates a new file material.
func NewFileMaterial(
	lessonID uuid.UUID,
	fileName string,
	fileURL *string,
	fileSize int64,
	fileType *string,
	sortOrder int64,
) *FileMaterial {
	return &FileMaterial{
		BaseModel: BaseModel{
			ID: uuid.New(),
		},

		LessonID: lessonID,

		FileName:  fileName,
		FileURL:   fileURL,
		FileSize:  fileSize,
		FileType:  fileType,
		SortOrder: sortOrder,
	}
}

