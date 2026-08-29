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

// ============================================================
// UPDATE PROGRESS
// ============================================================

func (s *lessonProgressService) UpdateProgress(
	userID uuid.UUID,
	lessonID uuid.UUID,
	req progress.UpdateProgressRequest,
) (*progress.ProgressResponse, error) {

	if req.WatchedSeconds < 0 {
		return nil, errors.New("watched_seconds cannot be negative")
	}

	// --------------------------------------------------------
	// 1. Find lesson
	// --------------------------------------------------------

	lesson, err := s.lessonRepository.FindByID(lessonID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("lesson not found")
		}

		return nil, err
	}

	// --------------------------------------------------------
	// 2. Check enrollment
	// --------------------------------------------------------

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

	// --------------------------------------------------------
	// 3. Find existing progress
	// --------------------------------------------------------

	existing, err := s.progressRepository.FindByEnrollmentAndLesson(
		enrollment.ID,
		lessonID,
	)

	// --------------------------------------------------------
	// 4. Create new progress
	// --------------------------------------------------------

	if errors.Is(err, pgx.ErrNoRows) {

		newProgress := &models.LessonProgress{
			BaseModel: models.BaseModel{
				ID: uuid.New(),
			},

			EnrollmentID: enrollment.ID,
			LessonID:     lessonID,

			Completed:      req.Completed,
			WatchedSeconds: int64(req.WatchedSeconds),
		}

		if err := s.progressRepository.Create(newProgress); err != nil {
			return nil, err
		}

		return mapProgressToResponse(newProgress), nil
	}

	// --------------------------------------------------------
	// 5. Other database error
	// --------------------------------------------------------

	if err != nil {
		return nil, err
	}

	// --------------------------------------------------------
	// 6. Update existing progress
	// --------------------------------------------------------

	existing.WatchedSeconds = int64(req.WatchedSeconds)
	existing.Completed = req.Completed

	if err := s.progressRepository.Update(existing); err != nil {
		return nil, err
	}

	return mapProgressToResponse(existing), nil
}

// ============================================================
// GET LESSON PROGRESS
// ============================================================

func (s *lessonProgressService) GetLessonProgress(
	userID uuid.UUID,
	lessonID uuid.UUID,
) (*progress.ProgressResponse, error) {

	// --------------------------------------------------------
	// 1. Find lesson
	// --------------------------------------------------------

	lesson, err := s.lessonRepository.FindByID(lessonID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("lesson not found")
		}

		return nil, err
	}

	// --------------------------------------------------------
	// 2. Check enrollment
	// --------------------------------------------------------

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

	// --------------------------------------------------------
	// 3. Find progress
	// --------------------------------------------------------

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

// ============================================================
// GET COURSE PROGRESS
// ============================================================

func (s *lessonProgressService) GetCourseProgress(
	userID uuid.UUID,
	courseID uuid.UUID,
) (*progress.CourseProgressResponse, error) {

	// --------------------------------------------------------
	// 1. Check course
	// --------------------------------------------------------

	courseModel, err := s.courseRepository.FindByID(courseID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("course not found")
		}

		return nil, err
	}

	// --------------------------------------------------------
	// 2. Check enrollment
	// --------------------------------------------------------

	enrollment, err := s.enrollmentRepository.FindByUserAndCourse(
		userID,
		courseID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("you are not enrolled in this course")
		}

		return nil, err
	}

	// --------------------------------------------------------
	// 3. Count total lessons
	// --------------------------------------------------------

	totalLessons := 0

	for _, section := range courseModel.Sections {
		totalLessons += len(section.Lessons)
	}

	// --------------------------------------------------------
	// 4. Get progress records
	// --------------------------------------------------------

	progresses, err := s.progressRepository.FindByEnrollment(
		enrollment.ID,
	)

	if err != nil {
		return nil, err
	}

	// --------------------------------------------------------
	// 5. Count completed lessons
	// --------------------------------------------------------

	completedLessons := 0

	for _, item := range progresses {
		if item.Completed {
			completedLessons++
		}
	}

	// --------------------------------------------------------
	// 6. Calculate percentage
	// --------------------------------------------------------

	var percent float64

	if totalLessons > 0 {
		percent = float64(completedLessons) /
			float64(totalLessons) *
			100
	}

	// --------------------------------------------------------
	// 7. Response
	// --------------------------------------------------------

	return &progress.CourseProgressResponse{
		CourseID:         courseID.String(),
		TotalLessons:     totalLessons,
		CompletedLessons: completedLessons,
		ProgressPercent:  percent,
	}, nil
}

// ============================================================
// MAP MODEL -> RESPONSE
// ============================================================

func mapProgressToResponse(
	item *models.LessonProgress,
) *progress.ProgressResponse {

	return &progress.ProgressResponse{
		ID:             item.ID.String(),
		EnrollmentID:   item.EnrollmentID.String(),
		LessonID:       item.LessonID.String(),
		Completed:      item.Completed,
		WatchedSeconds: int(item.WatchedSeconds),
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}
}
