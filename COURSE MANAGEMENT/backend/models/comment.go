package models

import "github.com/google/uuid"

type Comment struct {
	BaseModel

	UserID   uuid.UUID `json:"user_id"`
	CourseID uuid.UUID `json:"course_id"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"`

	User   User    `json:"user,omitempty"`
	Course Course  `json:"course,omitempty"`
	Parent *Comment `json:"parent,omitempty"`

	Content string `json:"content"`
}