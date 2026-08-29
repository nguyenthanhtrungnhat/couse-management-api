package models

import "github.com/google/uuid"

type Review struct {
	BaseModel

	UserID   uuid.UUID `json:"user_id"`
	CourseID uuid.UUID `json:"course_id"`

	User   User   `json:"user,omitempty"`
	Course Course `json:"course,omitempty"`

	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}