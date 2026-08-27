package services

import (
	"errors"
	"strings"

	"course-management/dto/section"
	"course-management/models"
	"course-management/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseSectionService interface {
	CreateSection(
		userID uuid.UUID,
		req section.CreateSectionRequest,
	) (*section.SectionResponse, error)

	GetSection(
		id uuid.UUID,
	) (*section.SectionResponse, error)

	GetSectionsByCourse(
		courseID uuid.UUID,
	) ([]section.SectionResponse, error)

	UpdateSection(
		userID uuid.UUID,
		id uuid.UUID,
		req section.UpdateSectionRequest,
	) (*section.SectionResponse, error)

	DeleteSection(
		userID uuid.UUID,
		id uuid.UUID,
	) error
}

type courseSectionService struct {
	sectionRepository repositories.CourseSectionRepository
	courseRepository   repositories.CourseRepository
}

func NewCourseSectionService(
	sectionRepository repositories.CourseSectionRepository,
	courseRepository repositories.CourseRepository,
) CourseSectionService {

	return &courseSectionService{
		sectionRepository: sectionRepository,
		courseRepository:   courseRepository,
	}
}

func (s *courseSectionService) CreateSection(
	userID uuid.UUID,
	req section.CreateSectionRequest,
) (*section.SectionResponse, error) {

	courseID, err := uuid.Parse(req.CourseID)

	if err != nil {
		return nil, errors.New("invalid course id")
	}

	title := strings.TrimSpace(req.Title)

	if title == "" {
		return nil, errors.New("section title is required")
	}

	courseModel, err := s.courseRepository.FindByID(courseID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}

		return nil, err
	}

	if courseModel.InstructorID != userID {
		return nil, errors.New("you are not the instructor of this course")
	}

	newSection := &models.CourseSection{
		CourseID:  courseID,
		Title:     title,
		SortOrder: req.SortOrder,
	}

	if err := s.sectionRepository.Create(newSection); err != nil {
		return nil, err
	}

	return mapSectionToResponse(newSection), nil
}

func (s *courseSectionService) GetSection(
	id uuid.UUID,
) (*section.SectionResponse, error) {

	result, err := s.sectionRepository.FindByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("section not found")
		}

		return nil, err
	}

	return mapSectionToResponse(result), nil
}

func (s *courseSectionService) GetSectionsByCourse(
	courseID uuid.UUID,
) ([]section.SectionResponse, error) {

	sections, err := s.sectionRepository.FindByCourseID(courseID)

	if err != nil {
		return nil, err
	}

	result := make(
		[]section.SectionResponse,
		0,
		len(sections),
	)

	for i := range sections {
		result = append(
			result,
			*mapSectionToResponse(&sections[i]),
		)
	}

	return result, nil
}

func (s *courseSectionService) UpdateSection(
	userID uuid.UUID,
	id uuid.UUID,
	req section.UpdateSectionRequest,
) (*section.SectionResponse, error) {

	existing, err := s.sectionRepository.FindByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("section not found")
		}

		return nil, err
	}

	courseModel, err := s.courseRepository.FindByID(
		existing.CourseID,
	)

	if err != nil {
		return nil, err
	}

	if courseModel.InstructorID != userID {
		return nil, errors.New(
			"you are not the instructor of this course",
		)
	}

	title := strings.TrimSpace(req.Title)

	if title == "" {
		return nil, errors.New("section title is required")
	}

	existing.Title = title
	existing.SortOrder = req.SortOrder

	if err := s.sectionRepository.Update(existing); err != nil {
		return nil, err
	}

	return mapSectionToResponse(existing), nil
}

func (s *courseSectionService) DeleteSection(
	userID uuid.UUID,
	id uuid.UUID,
) error {

	existing, err := s.sectionRepository.FindByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("section not found")
		}

		return err
	}

	courseModel, err := s.courseRepository.FindByID(
		existing.CourseID,
	)

	if err != nil {
		return err
	}

	if courseModel.InstructorID != userID {
		return errors.New(
			"you are not the instructor of this course",
		)
	}

	return s.sectionRepository.Delete(id)
}

func mapSectionToResponse(
	model *models.CourseSection,
) *section.SectionResponse {

	return &section.SectionResponse{
		ID:        model.ID.String(),
		CourseID:  model.CourseID.String(),
		Title:     model.Title,
		SortOrder: model.SortOrder,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}