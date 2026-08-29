package models

import (
	"github.com/google/uuid"
)

type Course struct {
	BaseModel

	InstructorID uuid.UUID `json:"instructor_id"`
	CategoryID   uuid.UUID `json:"category_id"`

	Title           string  `json:"title"`
	Slug            string  `json:"slug"`
	Description     *string `json:"description,omitempty"`
	ThumbnailURL    *string `json:"thumbnail_url,omitempty"`
	PreviewVideoURL *string `json:"preview_video_url,omitempty"`

	Price         float64 `json:"price"`
	Status        string  `json:"status"`
	AverageRating float64 `json:"average_rating"`
	TotalStudents int64   `json:"total_students"`

	Instructor *User     `json:"instructor,omitempty"`
	Category   *Category `json:"category,omitempty"`

	Sections []CourseSection `json:"sections,omitempty"`
	Reviews  []Review        `json:"reviews,omitempty"`
	Comments []Comment       `json:"comments,omitempty"`
}

// TableName returns the database table name.
func (Course) TableName() string {
	return "courses"
}

// NewCourse creates a new Course model.
func NewCourse(
	instructorID uuid.UUID,
	categoryID uuid.UUID,
	title string,
	slug string,
	description *string,
	thumbnailURL *string,
	previewVideoURL *string,
	price float64,
	status string,
) *Course {
	return &Course{
		BaseModel: BaseModel{
			ID: uuid.New(),
		},

		InstructorID: instructorID,
		CategoryID:   categoryID,

		Title:           title,
		Slug:            slug,
		Description:     description,
		ThumbnailURL:    thumbnailURL,
		PreviewVideoURL: previewVideoURL,

		Price:  price,
		Status: status,

		AverageRating: 0,
		TotalStudents: 0,
	}
}
