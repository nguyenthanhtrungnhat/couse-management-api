package models

import "github.com/google/uuid"

type Comment struct {
	BaseModel

	UserID uuid.UUID `gorm:"type:uuid;index;not null"`

	CourseID uuid.UUID `gorm:"type:uuid;index;not null"`

	ParentID *uuid.UUID `gorm:"type:uuid;index"`

	User User

	Course Course

	Parent *Comment `gorm:"foreignKey:ParentID"`

	Content string

	Replies []Comment `gorm:"foreignKey:ParentID"`
}