
package models

import "github.com/google/uuid"

type Enrollment struct {
	BaseModel

	UserID   uuid.UUID `json:"user_id"`
	CourseID uuid.UUID `json:"course_id"`

	User   *User   `json:"user,omitempty"`
	Course *Course `json:"course,omitempty"`
}

// TableName returns the database table name.
func (Enrollment) TableName() string {
	return "enrollments"
}

// NewEnrollment creates a new enrollment.
func NewEnrollment(
	userID uuid.UUID,
	courseID uuid.UUID,
) *Enrollment {
	return &Enrollment{
		BaseModel: BaseModel{
			ID: uuid.New(),
		},

		UserID:   userID,
		CourseID: courseID,
	}
}

