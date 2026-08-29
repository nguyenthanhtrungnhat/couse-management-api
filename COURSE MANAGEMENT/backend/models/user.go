
package models

import "github.com/google/uuid"

type User struct {
	BaseModel

	RoleID uuid.UUID `json:"role_id"`

	FullName     string  `json:"full_name"`
	Email        string  `json:"email"`
	PasswordHash *string `json:"-"`
	AvatarURL    *string `json:"avatar_url,omitempty"`
	Provider     string  `json:"provider"`

	Role Role `json:"role,omitempty"`
}

// TableName returns the database table name.
func (User) TableName() string {
	return "users"
}

// NewUser creates a new User model.
func NewUser(
	roleID uuid.UUID,
	fullName string,
	email string,
	passwordHash *string,
	provider string,
) *User {
	return &User{
		BaseModel: BaseModel{
			ID: uuid.New(),
		},
		RoleID:       roleID,
		FullName:     fullName,
		Email:        email,
		PasswordHash: passwordHash,
		Provider:     provider,
	}
}

