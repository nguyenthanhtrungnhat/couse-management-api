package models

import "github.com/google/uuid"

type Review struct {
	BaseModel

	UserID   uuid.UUID `json:"user_id"`
	CourseID uuid.UUID `json:"course_id"`

	Rating  int64  `json:"rating"`
	Comment string `json:"comment"`

	User   *User   `json:"user,omitempty"`
	Course *Course `json:"course,omitempty"`
}

// TableName returns the database table name.
func (Review) TableName() string {
	return "reviews"
}

// NewReview creates a new review.
func NewReview(
	userID uuid.UUID,
	courseID uuid.UUID,
	rating int64,
	comment string,
) *Review {
	return &Review{
		BaseModel: BaseModel{
			ID: uuid.New(),
		},

		UserID:   userID,
		CourseID: courseID,
		Rating:   rating,
		Comment:  comment,
	}
}
