package section

import (
	"time"

	"course-management/dto/lesson"
)

type SectionResponse struct {
	ID        string    `json:"id"`
	CourseID  string    `json:"course_id"`
	Title     string    `json:"title"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Lessons []lesson.LessonResponse `json:"lessons"`
}
