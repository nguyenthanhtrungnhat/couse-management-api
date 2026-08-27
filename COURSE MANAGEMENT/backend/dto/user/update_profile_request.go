package user

type UpdateProfileRequest struct {
	FullName  string  `json:"full_name" validate:"required,min=3,max=255"`
	AvatarURL *string `json:"avatar_url"`
}