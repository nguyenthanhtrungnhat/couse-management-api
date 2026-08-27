package repositories

import (
	"course-management/config"
	"course-management/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LessonProgressRepository interface {
	FindByEnrollmentAndLesson(
		enrollmentID uuid.UUID,
		lessonID uuid.UUID,
	) (*models.LessonProgress, error)

	Create(progress *models.LessonProgress) error

	Update(progress *models.LessonProgress) error

	FindByEnrollment(
		enrollmentID uuid.UUID,
	) ([]models.LessonProgress, error)
}

type lessonProgressRepository struct {
	db *gorm.DB
}

func NewLessonProgressRepository() LessonProgressRepository {
	return &lessonProgressRepository{
		db: config.DB,
	}
}

func (r *lessonProgressRepository) FindByEnrollmentAndLesson(
	enrollmentID uuid.UUID,
	lessonID uuid.UUID,
) (*models.LessonProgress, error) {

	var progress models.LessonProgress

	err := r.db.
		Where(
			"enrollment_id = ? AND lesson_id = ?",
			enrollmentID,
			lessonID,
		).
		First(&progress).
		Error

	if err != nil {
		return nil, err
	}

	return &progress, nil
}

func (r *lessonProgressRepository) Create(
	progress *models.LessonProgress,
) error {
	return r.db.Create(progress).Error
}

func (r *lessonProgressRepository) Update(
	progress *models.LessonProgress,
) error {
	return r.db.Save(progress).Error
}

func (r *lessonProgressRepository) FindByEnrollment(
	enrollmentID uuid.UUID,
) ([]models.LessonProgress, error) {

	var progresses []models.LessonProgress

	err := r.db.
		Where("enrollment_id = ?", enrollmentID).
		Order("created_at ASC").
		Find(&progresses).
		Error

	if err != nil {
		return nil, err
	}

	return progresses, nil
}