package services

import (
	"errors"

	"course-management/dto/progress"
	"course-management/models"
	"course-management/repositories"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type LessonProgressService interface {
	UpdateProgress(
		userID uuid.UUID,
		lessonID uuid.UUID,
		req progress.UpdateProgressRequest,
	) (*progress.ProgressResponse, error)

	GetLessonProgress(
		userID uuid.UUID,
		lessonID uuid.UUID,
	) (*progress.ProgressResponse, error)

	GetCourseProgress(
		userID uuid.UUID,
		courseID uuid.UUID,
	) (*progress.CourseProgressResponse, error)
}

type lessonProgressService struct {
	progressRepository   repositories.LessonProgressRepository
	enrollmentRepository repositories.EnrollmentRepository
	lessonRepository     repositories.LessonRepository
	courseRepository     repositories.CourseRepository
}

func NewLessonProgressService(
	progressRepository repositories.LessonProgressRepository,
	enrollmentRepository repositories.EnrollmentRepository,
	lessonRepository repositories.LessonRepository,
	courseRepository repositories.CourseRepository,
) LessonProgressService {
	return &lessonProgressService{
		progressRepository:   progressRepository,
		enrollmentRepository: enrollmentRepository,
		lessonRepository:     lessonRepository,
		courseRepository:     courseRepository,
	}
}

func (s *lessonProgressService) UpdateProgress(
	userID uuid.UUID,
	lessonID uuid.UUID,
	req progress.UpdateProgressRequest,
) (*progress.ProgressResponse, error) {

	if req.WatchedSeconds < 0 {
		return nil, errors.New("watched_seconds cannot be negative")
	}

	lesson, err := s.lessonRepository.FindByID(lessonID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("lesson not found")
		}
		return nil, err
	}

	enrollment, err := s.enrollmentRepository.FindByUserAndCourse(
		userID,
		lesson.Section.CourseID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("you are not enrolled in this course")
		}
		return nil, err
	}

	existing, err := s.progressRepository.FindByEnrollmentAndLesson(
		enrollment.ID,
		lessonID,
	)

	if errors.Is(err, pgx.ErrNoRows) {

		newProgress := &models.LessonProgress{
			EnrollmentID:   enrollment.ID,
			LessonID:       lessonID,
			WatchedSeconds: req.WatchedSeconds,
			Completed:      req.Completed,
		}

		if err := s.progressRepository.Create(newProgress); err != nil {
			return nil, err
		}

		return mapProgressToResponse(newProgress), nil
	}

	if err != nil {
		return nil, err
	}

	existing.WatchedSeconds = req.WatchedSeconds
	existing.Completed = req.Completed

	if err := s.progressRepository.Update(existing); err != nil {
		return nil, err
	}

	return mapProgressToResponse(existing), nil
}

func (s *lessonProgressService) GetLessonProgress(
	userID uuid.UUID,
	lessonID uuid.UUID,
) (*progress.ProgressResponse, error) {

	lesson, err := s.lessonRepository.FindByID(lessonID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("lesson not found")
		}
		return nil, err
	}

	enrollment, err := s.enrollmentRepository.FindByUserAndCourse(
		userID,
		lesson.Section.CourseID,
	)
	if err != nil {
		return nil, errors.New("you are not enrolled in this course")
	}

	item, err := s.progressRepository.FindByEnrollmentAndLesson(
		enrollment.ID,
		lessonID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("progress not found")
		}

		return nil, err
	}

	return mapProgressToResponse(item), nil
}

func (s *lessonProgressService) GetCourseProgress(
	userID uuid.UUID,
	courseID uuid.UUID,
) (*progress.CourseProgressResponse, error) {

	_, err := s.courseRepository.FindByID(courseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}

	enrollment, err := s.enrollmentRepository.FindByUserAndCourse(
		userID,
		courseID,
	)
	if err != nil {
		return nil, errors.New("you are not enrolled in this course")
	}

	courseModel, err := s.courseRepository.FindByID(courseID)
	if err != nil {
		return nil, err
	}

	totalLessons := 0

	for _, section := range courseModel.Sections {
		totalLessons += len(section.Lessons)
	}

	progresses, err := s.progressRepository.FindByEnrollment(
		enrollment.ID,
	)
	if err != nil {
		return nil, err
	}

	completedLessons := 0

	for _, item := range progresses {
		if item.Completed {
			completedLessons++
		}
	}

	var percent float64

	if totalLessons > 0 {
		percent = float64(completedLessons) /
			float64(totalLessons) *
			100
	}

	return &progress.CourseProgressResponse{
		CourseID:         courseID.String(),
		TotalLessons:     totalLessons,
		CompletedLessons: completedLessons,
		ProgressPercent:  percent,
	}, nil
}

func mapProgressToResponse(
	item *models.LessonProgress,
) *progress.ProgressResponse {

	return &progress.ProgressResponse{
		ID:             item.ID.String(),
		EnrollmentID:   item.EnrollmentID.String(),
		LessonID:       item.LessonID.String(),
		Completed:      item.Completed,
		WatchedSeconds: item.WatchedSeconds,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}
}
