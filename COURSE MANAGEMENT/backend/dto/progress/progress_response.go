package progress

import "time"

type ProgressResponse struct {
	ID             string    `json:"id"`
	EnrollmentID   string    `json:"enrollment_id"`
	LessonID       string    `json:"lesson_id"`
	Completed      bool      `json:"completed"`
	WatchedSeconds int       `json:"watched_seconds"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CourseProgressResponse struct {
	CourseID         string  `json:"course_id"`
	TotalLessons     int     `json:"total_lessons"`
	CompletedLessons int     `json:"completed_lessons"`
	ProgressPercent  float64 `json:"progress_percent"`
}
