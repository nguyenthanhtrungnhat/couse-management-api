package repositories

import (
	"course-management/config"
	"course-management/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EnrollmentRepository interface {
	Create(enrollment *models.Enrollment) error
	FindByID(id uuid.UUID) (*models.Enrollment, error)
	FindByUserAndCourse(userID uuid.UUID, courseID uuid.UUID) (*models.Enrollment, error)
	FindByUser(userID uuid.UUID) ([]models.Enrollment, error)
	FindByCourse(courseID uuid.UUID) ([]models.Enrollment, error)
	Exists(userID uuid.UUID, courseID uuid.UUID) (bool, error)
	Delete(id uuid.UUID) error
	CountByCourse(courseID uuid.UUID) (int64, error)
}

type enrollmentRepository struct {
	db *gorm.DB
}

func NewEnrollmentRepository() EnrollmentRepository {
	return &enrollmentRepository{
		db: config.DB,
	}
}

func (r *enrollmentRepository) Create(
	enrollment *models.Enrollment,
) error {
	return r.db.Create(enrollment).Error
}

func (r *enrollmentRepository) FindByID(
	id uuid.UUID,
) (*models.Enrollment, error) {

	var enrollment models.Enrollment

	err := r.db.
		Preload("User").
		Preload("Course").
		First(&enrollment, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return &enrollment, nil
}

func (r *enrollmentRepository) FindByUserAndCourse(
	userID uuid.UUID,
	courseID uuid.UUID,
) (*models.Enrollment, error) {

	var enrollment models.Enrollment

	err := r.db.
		Preload("Course").
		Where(
			"user_id = ? AND course_id = ?",
			userID,
			courseID,
		).
		First(&enrollment).
		Error

	if err != nil {
		return nil, err
	}

	return &enrollment, nil
}

func (r *enrollmentRepository) FindByUser(
	userID uuid.UUID,
) ([]models.Enrollment, error) {

	var enrollments []models.Enrollment

	err := r.db.
		Preload("Course").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&enrollments).
		Error

	return enrollments, err
}

func (r *enrollmentRepository) FindByCourse(
	courseID uuid.UUID,
) ([]models.Enrollment, error) {

	var enrollments []models.Enrollment

	err := r.db.
		Preload("User").
		Where("course_id = ?", courseID).
		Order("created_at DESC").
		Find(&enrollments).
		Error

	return enrollments, err
}

func (r *enrollmentRepository) Exists(
	userID uuid.UUID,
	courseID uuid.UUID,
) (bool, error) {

	var count int64

	err := r.db.
		Model(&models.Enrollment{}).
		Where(
			"user_id = ? AND course_id = ?",
			userID,
			courseID,
		).
		Count(&count).
		Error

	return count > 0, err
}

func (r *enrollmentRepository) Delete(
	id uuid.UUID,
) error {
	return r.db.
		Delete(&models.Enrollment{}, "id = ?", id).
		Error
}

func (r *enrollmentRepository) CountByCourse(
	courseID uuid.UUID,
) (int64, error) {

	var count int64

	err := r.db.
		Model(&models.Enrollment{}).
		Where("course_id = ?", courseID).
		Count(&count).
		Error

	return count, err
}