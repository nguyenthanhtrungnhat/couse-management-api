package services

import (
	authDTO "course-management/dto/auth"
	"course-management/models"
	"course-management/repositories"
	"course-management/utils"
	"course-management/constants"
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

// Register tạo tài khoản mới
func (s *authService) Register(req authDTO.RegisterRequest) (*authDTO.AuthResponse, error) {

	// Kiểm tra email đã tồn tại
	_, err := s.userRepo.FindByEmail(req.Email)

	if err == nil {
		return nil, constants.ErrEmailExists
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Lấy role student
	role, err := s.roleRepo.FindByName("student")
	if err != nil {
		return nil, err
	}

	// Hash password
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Tạo user
	user := models.User{
		RoleID:       role.ID,
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: &hash,
		Provider:     "local",
	}

	// Lưu database
	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	// Load Role để trả response
	user.Role = *role

	// Tạo JWT
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

// Login
func (s *authService) Login(req authDTO.LoginRequest) (*authDTO.AuthResponse, error) {

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

	// Kiểm tra mật khẩu
	if !utils.CheckPassword(*user.PasswordHash, req.Password) {
		return nil, constants.ErrInvalidCredential
	}

	// Sinh JWT
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