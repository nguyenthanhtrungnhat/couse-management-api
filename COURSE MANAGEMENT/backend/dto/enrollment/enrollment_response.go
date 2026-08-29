package enrollment

import "time"

type EnrollmentResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	CourseID    string    `json:"course_id"`
	CourseTitle string    `json:"course_title,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
