package models

import "github.com/google/uuid"

type User struct {
	BaseModel

	RoleID uuid.UUID `json:"role_id"`
	Role   Role      `json:"role,omitempty"`

	FullName     string  `json:"full_name"`
	Email        string  `json:"email"`
	PasswordHash *string `json:"-"`
	AvatarURL    *string `json:"avatar_url,omitempty"`
	Provider     string  `json:"provider"`

	Courses []Course `json:"courses,omitempty"`
}