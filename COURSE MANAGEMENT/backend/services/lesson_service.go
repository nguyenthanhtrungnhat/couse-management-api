package services

import (
	"errors"
	"strings"

	"course-management/dto/lesson"
	"course-management/models"
	"course-management/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LessonService interface {
	CreateLesson(
		userID uuid.UUID,
		req lesson.CreateLessonRequest,
	) (*lesson.LessonResponse, error)

	GetLesson(
		id uuid.UUID,
	) (*lesson.LessonResponse, error)

	GetLessonsBySection(
		sectionID uuid.UUID,
	) ([]lesson.LessonResponse, error)

	UpdateLesson(
		userID uuid.UUID,
		id uuid.UUID,
		req lesson.UpdateLessonRequest,
	) (*lesson.LessonResponse, error)

	DeleteLesson(
		userID uuid.UUID,
		id uuid.UUID,
	) error
}

type lessonService struct {
	lessonRepository  repositories.LessonRepository
	sectionRepository repositories.CourseSectionRepository
	courseRepository  repositories.CourseRepository
}

func NewLessonService(
	lessonRepository repositories.LessonRepository,
	sectionRepository repositories.CourseSectionRepository,
	courseRepository repositories.CourseRepository,
) LessonService {

	return &lessonService{
		lessonRepository:  lessonRepository,
		sectionRepository: sectionRepository,
		courseRepository:  courseRepository,
	}
}

func (s *lessonService) CreateLesson(
	userID uuid.UUID,
	req lesson.CreateLessonRequest,
) (*lesson.LessonResponse, error) {

	sectionID, err := uuid.Parse(req.SectionID)

	if err != nil {
		return nil, errors.New("invalid section id")
	}

	title := strings.TrimSpace(req.Title)
	videoURL := strings.TrimSpace(req.VideoURL)

	if title == "" {
		return nil, errors.New("lesson title is required")
	}

	if videoURL == "" {
		return nil, errors.New("video URL is required")
	}

	section, err := s.sectionRepository.FindByID(sectionID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("section not found")
		}

		return nil, err
	}

	course, err := s.courseRepository.FindByID(
		section.CourseID,
	)

	if err != nil {
		return nil, err
	}

	if course.InstructorID != userID {
		return nil, errors.New(
			"you are not the instructor of this course",
		)
	}

	newLesson := &models.Lesson{
		SectionID:       sectionID,
		Title:           title,
		VideoURL:        videoURL,
		DurationSeconds: req.DurationSeconds,
		IsPreview:       req.IsPreview,
		SortOrder:       req.SortOrder,
	}

	if err := s.lessonRepository.Create(newLesson); err != nil {
		return nil, err
	}

	return mapLessonToResponse(newLesson), nil
}

func (s *lessonService) GetLesson(
	id uuid.UUID,
) (*lesson.LessonResponse, error) {

	result, err := s.lessonRepository.FindByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lesson not found")
		}

		return nil, err
	}

	return mapLessonToResponse(result), nil
}

func (s *lessonService) GetLessonsBySection(
	sectionID uuid.UUID,
) ([]lesson.LessonResponse, error) {

	lessons, err := s.lessonRepository.FindBySectionID(
		sectionID,
	)

	if err != nil {
		return nil, err
	}

	result := make(
		[]lesson.LessonResponse,
		0,
		len(lessons),
	)

	for i := range lessons {
		result = append(
			result,
			*mapLessonToResponse(&lessons[i]),
		)
	}

	return result, nil
}

func (s *lessonService) UpdateLesson(
	userID uuid.UUID,
	id uuid.UUID,
	req lesson.UpdateLessonRequest,
) (*lesson.LessonResponse, error) {

	existing, err := s.lessonRepository.FindByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lesson not found")
		}

		return nil, err
	}

	section, err := s.sectionRepository.FindByID(
		existing.SectionID,
	)

	if err != nil {
		return nil, err
	}

	course, err := s.courseRepository.FindByID(
		section.CourseID,
	)

	if err != nil {
		return nil, err
	}

	if course.InstructorID != userID {
		return nil, errors.New(
			"you are not the instructor of this course",
		)
	}

	title := strings.TrimSpace(req.Title)
	videoURL := strings.TrimSpace(req.VideoURL)

	if title == "" {
		return nil, errors.New("lesson title is required")
	}

	if videoURL == "" {
		return nil, errors.New("video URL is required")
	}

	existing.Title = title
	existing.VideoURL = videoURL
	existing.DurationSeconds = req.DurationSeconds
	existing.IsPreview = req.IsPreview
	existing.SortOrder = req.SortOrder

	if err := s.lessonRepository.Update(existing); err != nil {
		return nil, err
	}

	return mapLessonToResponse(existing), nil
}

func (s *lessonService) DeleteLesson(
	userID uuid.UUID,
	id uuid.UUID,
) error {

	existing, err := s.lessonRepository.FindByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("lesson not found")
		}

		return err
	}

	section, err := s.sectionRepository.FindByID(
		existing.SectionID,
	)

	if err != nil {
		return err
	}

	course, err := s.courseRepository.FindByID(
		section.CourseID,
	)

	if err != nil {
		return err
	}

	if course.InstructorID != userID {
		return errors.New(
			"you are not the instructor of this course",
		)
	}

	return s.lessonRepository.Delete(id)
}

func mapLessonToResponse(
	model *models.Lesson,
) *lesson.LessonResponse {

	return &lesson.LessonResponse{
		ID:              model.ID.String(),
		SectionID:       model.SectionID.String(),
		Title:           model.Title,
		VideoURL:        model.VideoURL,
		DurationSeconds: model.DurationSeconds,
		IsPreview:       model.IsPreview,
		SortOrder:       model.SortOrder,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}
}