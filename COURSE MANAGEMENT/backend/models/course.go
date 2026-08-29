package models

import "github.com/google/uuid"

type Course struct {
	BaseModel

	InstructorID uuid.UUID `json:"instructor_id"`
	Instructor   User      `json:"instructor,omitempty"`

	CategoryID uuid.UUID `json:"category_id"`
	Category   Category  `json:"category,omitempty"`

	Title           string  `json:"title"`
	Slug            string  `json:"slug"`
	Description     *string `json:"description"`
	ThumbnailURL    *string `json:"thumbnail_url"`
	PreviewVideoURL *string `json:"preview_video_url"`

	Price         float64 `json:"price"`
	Status        string  `json:"status"`
	AverageRating float64 `json:"average_rating"`
	TotalStudents int     `json:"total_students"`

	Sections []CourseSection `json:"sections,omitempty"`
	Reviews  []Review        `json:"reviews,omitempty"`
}