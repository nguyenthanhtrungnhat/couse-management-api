package services

import (
	userDTO "course-management/dto/user"
	"course-management/models"
	"course-management/repositories"

	"github.com/google/uuid"
)

type UserService interface {
	GetProfile(userID uuid.UUID) (*models.User, error)
	UpdateProfile(userID uuid.UUID, req userDTO.UpdateProfileRequest) (*models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService() UserService {
	return &userService{
		userRepo: repositories.NewUserRepository(),
	}
}

func (s *userService) GetProfile(userID uuid.UUID) (*models.User, error) {

	user, err := s.userRepo.FindByID(userID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateProfile(
	userID uuid.UUID,
	req userDTO.UpdateProfileRequest,
) (*models.User, error) {

	user, err := s.userRepo.FindByID(userID)

	if err != nil {
		return nil, err
	}

	user.FullName = req.FullName

	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}