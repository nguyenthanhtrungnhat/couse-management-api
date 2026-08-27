package models

import "github.com/google/uuid"

type User struct {
	BaseModel

	RoleID uuid.UUID `gorm:"type:uuid;index;not null"`
	Role   Role

	FullName string `gorm:"size:255;not null"`

	Email string `gorm:"size:255;uniqueIndex;not null"`

	PasswordHash *string

	AvatarURL *string

	Provider string `gorm:"size:20;default:local"`

	Courses []Course `gorm:"foreignKey:InstructorID"`
}
