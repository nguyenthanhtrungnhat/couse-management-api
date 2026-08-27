package models

import "github.com/google/uuid"

type Review struct {
	BaseModel

	UserID uuid.UUID `gorm:"type:uuid;index;not null"`

	CourseID uuid.UUID `gorm:"type:uuid;index;not null"`

	User User

	Course Course

	Rating int

	Comment string
}