package models

import "github.com/google/uuid"

type Course struct {
	BaseModel

	InstructorID uuid.UUID `gorm:"type:uuid;index;not null" json:"instructor_id"`
	Instructor   User      `json:"instructor"`

	CategoryID uuid.UUID `gorm:"type:uuid;index;not null" json:"category_id"`
	Category   Category   `json:"category"`

	Title string `gorm:"size:255;not null" json:"title"`

	Slug string `gorm:"size:255;uniqueIndex;not null" json:"slug"`

	Description *string `json:"description"`

	ThumbnailURL *string `json:"thumbnail_url"`

	PreviewVideoURL *string `json:"preview_video_url"`

	Price float64 `gorm:"type:numeric(12,2);default:0" json:"price"`

	Status string `gorm:"size:20;default:draft" json:"status"`

	AverageRating float64 `gorm:"type:numeric(3,2);default:0" json:"average_rating"`

	TotalStudents int `gorm:"default:0" json:"total_students"`

	Sections []CourseSection `gorm:"foreignKey:CourseID" json:"sections"`

	Reviews []Review `gorm:"foreignKey:CourseID" json:"reviews"`
}