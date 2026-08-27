package services

import (
	"course-management/constants"
	authDTO "course-management/dto/auth"
	"course-management/models"
	"course-management/repositories"
	"course-management/utils"
	"errors"

	"gorm.io/gorm"
)

type AuthService interface {
	Register(req authDTO.RegisterRequest) (*authDTO.AuthResponse, error)
	Login(req authDTO.LoginRequest) (*authDTO.AuthResponse, error)
}

type authService struct {
	userRepo repositories.UserRepository
	roleRepo repositories.RoleRepository
}

func NewAuthService() AuthService {
	return &authService{
		userRepo: repositories.NewUserRepository(),
		roleRepo: repositories.NewRoleRepository(),
	}
}

// ==========================================
// REGISTER
// ==========================================

func (s *authService) Register(
	req authDTO.RegisterRequest,
) (*authDTO.AuthResponse, error) {

	// Check email exists
	_, err := s.userRepo.FindByEmail(req.Email)

	if err == nil {
		return nil, constants.ErrEmailExists
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Default role = student
	role, err := s.roleRepo.FindByName("student")

	if err != nil {
		return nil, err
	}

	// Hash password
	hash, err := utils.HashPassword(req.Password)

	if err != nil {
		return nil, err
	}

	// Create user
	user := models.User{
		RoleID:       role.ID,
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: &hash,
		Provider:     "local",
	}

	// Save database
	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	// Load role
	user.Role = *role

	// Generate JWT
	accessToken, err := utils.GenerateAccessToken(
		user.ID.String(),
		role.Name,
	)

	if err != nil {
		return nil, err
	}

	return &authDTO.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: "",
		User: authDTO.UserResponse{
			ID:        user.ID.String(),
			FullName:  user.FullName,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
			Role:      role.Name,
		},
	}, nil
}

// ==========================================
// LOGIN
// ==========================================

func (s *authService) Login(
	req authDTO.LoginRequest,
) (*authDTO.AuthResponse, error) {

	user, err := s.userRepo.FindByEmail(req.Email)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrInvalidCredential
		}

		return nil, err
	}

	// Google account
	if user.PasswordHash == nil {
		return nil, constants.ErrGoogleAccount
	}

	// Check password
	if !utils.CheckPassword(
		*user.PasswordHash,
		req.Password,
	) {
		return nil, constants.ErrInvalidCredential
	}

	// Generate JWT
	accessToken, err := utils.GenerateAccessToken(
		user.ID.String(),
		user.Role.Name,
	)

	if err != nil {
		return nil, err
	}

	return &authDTO.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: "",
		User: authDTO.UserResponse{
			ID:        user.ID.String(),
			FullName:  user.FullName,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
			Role:      user.Role.Name,
		},
	}, nil
}
