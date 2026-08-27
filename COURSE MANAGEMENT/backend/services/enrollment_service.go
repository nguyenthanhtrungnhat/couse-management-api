package services

import (
	"errors"

	"course-management/dto/enrollment"
	"course-management/models"
	"course-management/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EnrollmentService interface {
	Enroll(userID uuid.UUID, courseID uuid.UUID) (*enrollment.EnrollmentResponse, error)
	GetMyEnrollments(userID uuid.UUID) ([]enrollment.EnrollmentResponse, error)
	GetEnrollmentByID(userID uuid.UUID, enrollmentID uuid.UUID) (*enrollment.EnrollmentResponse, error)
	Unenroll(userID uuid.UUID, enrollmentID uuid.UUID) error
}

type enrollmentService struct {
	enrollmentRepository repositories.EnrollmentRepository
	courseRepository     repositories.CourseRepository
}

func NewEnrollmentService(
	enrollmentRepository repositories.EnrollmentRepository,
	courseRepository repositories.CourseRepository,
) EnrollmentService {
	return &enrollmentService{
		enrollmentRepository: enrollmentRepository,
		courseRepository:     courseRepository,
	}
}

func (s *enrollmentService) Enroll(
	userID uuid.UUID,
	courseID uuid.UUID,
) (*enrollment.EnrollmentResponse, error) {

	// Course must exist.
	courseModel, err := s.courseRepository.FindByID(courseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}

		return nil, err
	}

	// Only published courses can be enrolled.
	if courseModel.Status != "published" {
		return nil, errors.New("course is not available for enrollment")
	}

	// User cannot enroll twice.
	exists, err := s.enrollmentRepository.Exists(
		userID,
		courseID,
	)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("already enrolled in this course")
	}

	model := &models.Enrollment{
		UserID:   userID,
		CourseID: courseID,
	}

	if err := s.enrollmentRepository.Create(model); err != nil {
		return nil, err
	}

	return s.mapEnrollment(model, courseModel.Title), nil
}

func (s *enrollmentService) GetMyEnrollments(
	userID uuid.UUID,
) ([]enrollment.EnrollmentResponse, error) {

	items, err := s.enrollmentRepository.FindByUser(userID)
	if err != nil {
		return nil, err
	}

	responses := make(
		[]enrollment.EnrollmentResponse,
		0,
		len(items),
	)

	for _, item := range items {
		title := ""

		if item.Course.ID != uuid.Nil {
			title = item.Course.Title
		}

		responses = append(
			responses,
			*s.mapEnrollment(&item, title),
		)
	}

	return responses, nil
}

func (s *enrollmentService) GetEnrollmentByID(
	userID uuid.UUID,
	enrollmentID uuid.UUID,
) (*enrollment.EnrollmentResponse, error) {

	item, err := s.enrollmentRepository.FindByID(enrollmentID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("enrollment not found")
		}

		return nil, err
	}

	if item.UserID != userID {
		return nil, errors.New("permission denied")
	}

	title := ""

	if item.Course.ID != uuid.Nil {
		title = item.Course.Title
	}

	return s.mapEnrollment(item, title), nil
}

func (s *enrollmentService) Unenroll(
	userID uuid.UUID,
	enrollmentID uuid.UUID,
) error {

	item, err := s.enrollmentRepository.FindByID(enrollmentID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("enrollment not found")
		}

		return err
	}

	if item.UserID != userID {
		return errors.New("permission denied")
	}

	return s.enrollmentRepository.Delete(enrollmentID)
}

func (s *enrollmentService) mapEnrollment(
	item *models.Enrollment,
	courseTitle string,
) *enrollment.EnrollmentResponse {

	return &enrollment.EnrollmentResponse{
		ID:          item.ID.String(),
		UserID:      item.UserID.String(),
		CourseID:    item.CourseID.String(),
		CourseTitle: courseTitle,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}