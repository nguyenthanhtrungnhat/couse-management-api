package services

import (
	"errors"
	"strings"

	"course-management/dto/lesson"
	"course-management/models"
	"course-management/repositories"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

// ============================================================
// CREATE LESSON
// ============================================================

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

	sectionModel, err := s.sectionRepository.FindByID(sectionID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("section not found")
		}

		return nil, err
	}

	courseModel, err := s.courseRepository.FindByID(
		sectionModel.CourseID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("course not found")
		}

		return nil, err
	}

	if courseModel.InstructorID != userID {
		return nil, errors.New(
			"you are not the instructor of this course",
		)
	}

	newLesson := &models.Lesson{
		SectionID: sectionID,
		Title:     title,

		// PostgreSQL TEXT nullable -> *string
		VideoURL: &videoURL,

		// PostgreSQL BIGINT -> int64
		DurationSeconds: int64(req.DurationSeconds),

		IsPreview: req.IsPreview,

		// PostgreSQL BIGINT -> int64
		SortOrder: int64(req.SortOrder),
	}

	if err := s.lessonRepository.Create(newLesson); err != nil {
		return nil, err
	}

	return mapLessonToResponse(newLesson), nil

}

// ============================================================
// GET LESSON
// ============================================================

func (s *lessonService) GetLesson(
	id uuid.UUID,
) (*lesson.LessonResponse, error) {

	result, err := s.lessonRepository.FindByID(id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("lesson not found")
		}

		return nil, err
	}

	return mapLessonToResponse(result), nil

}

// ============================================================
// GET LESSONS BY SECTION
// ============================================================

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

// ============================================================
// UPDATE LESSON
// ============================================================

func (s *lessonService) UpdateLesson(
	userID uuid.UUID,
	id uuid.UUID,
	req lesson.UpdateLessonRequest,
) (*lesson.LessonResponse, error) {

	existing, err := s.lessonRepository.FindByID(id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("lesson not found")
		}

		return nil, err
	}

	sectionModel, err := s.sectionRepository.FindByID(
		existing.SectionID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("section not found")
		}

		return nil, err
	}

	courseModel, err := s.courseRepository.FindByID(
		sectionModel.CourseID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("course not found")
		}

		return nil, err
	}

	if courseModel.InstructorID != userID {
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

	// PostgreSQL TEXT nullable -> *string
	existing.VideoURL = &videoURL

	// PostgreSQL BIGINT -> int64
	existing.DurationSeconds = int64(req.DurationSeconds)

	existing.IsPreview = req.IsPreview

	// PostgreSQL BIGINT -> int64
	existing.SortOrder = int64(req.SortOrder)

	if err := s.lessonRepository.Update(existing); err != nil {
		return nil, err
	}

	return mapLessonToResponse(existing), nil

}

// ============================================================
// DELETE LESSON
// ============================================================

func (s *lessonService) DeleteLesson(
	userID uuid.UUID,
	id uuid.UUID,
) error {

	existing, err := s.lessonRepository.FindByID(id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("lesson not found")
		}

		return err
	}

	sectionModel, err := s.sectionRepository.FindByID(
		existing.SectionID,
	)

	if err != nil {
		return err
	}

	courseModel, err := s.courseRepository.FindByID(
		sectionModel.CourseID,
	)

	if err != nil {
		return err
	}

	if courseModel.InstructorID != userID {
		return errors.New(
			"you are not the instructor of this course",
		)
	}

	return s.lessonRepository.Delete(id)

}

// ============================================================
// MAPPER
// ============================================================

func mapLessonToResponse(
	model *models.Lesson,
) *lesson.LessonResponse {

	videoURL := ""

	if model.VideoURL != nil {
		videoURL = *model.VideoURL
	}

	return &lesson.LessonResponse{
		ID:        model.ID.String(),
		SectionID: model.SectionID.String(),
		Title:     model.Title,

		VideoURL: videoURL,

		DurationSeconds: int(model.DurationSeconds),

		IsPreview: model.IsPreview,

		SortOrder: int(model.SortOrder),

		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}

}
