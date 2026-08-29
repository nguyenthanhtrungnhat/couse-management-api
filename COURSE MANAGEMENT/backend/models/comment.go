package models

import "github.com/google/uuid"

type Comment struct {
	BaseModel

	UserID   uuid.UUID  `json:"user_id"`
	CourseID uuid.UUID  `json:"course_id"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"`

	Content string `json:"content"`

	User   *User   `json:"user,omitempty"`
	Course *Course `json:"course,omitempty"`

	Parent  *Comment  `json:"parent,omitempty"`
	Replies []Comment `json:"replies,omitempty"`
}

// TableName returns the database table name.
func (Comment) TableName() string {
	return "comments"
}

// NewComment creates a new comment.
func NewComment(
	userID uuid.UUID,
	courseID uuid.UUID,
	parentID *uuid.UUID,
	content string,
) *Comment {
	return &Comment{
		BaseModel: BaseModel{
			ID: uuid.New(),
		},

		UserID:   userID,
		CourseID: courseID,
		ParentID: parentID,
		Content:  content,
	}
}
