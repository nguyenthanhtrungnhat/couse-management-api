package user

import "time"

type UserResponse struct {
	ID        string     `json:"id"`
	FullName  string     `json:"full_name"`
	Email     string     `json:"email"`
	AvatarURL *string    `json:"avatar_url,omitempty"`
	Provider  string     `json:"provider"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
}