package course

import "time"

type CourseResponse struct {
	ID              string     `json:"id"`
	InstructorID    string     `json:"instructor_id"`
	CategoryID      string     `json:"category_id"`
	Title           string     `json:"title"`
	Slug            string     `json:"slug"`
	Description     *string    `json:"description,omitempty"`
	ThumbnailURL    *string    `json:"thumbnail_url,omitempty"`
	PreviewVideoURL *string    `json:"preview_video_url,omitempty"`
	Price           float64    `json:"price"`
	Status          string     `json:"status"`
	AverageRating   float64    `json:"average_rating"`
	TotalStudents   int        `json:"total_students"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}