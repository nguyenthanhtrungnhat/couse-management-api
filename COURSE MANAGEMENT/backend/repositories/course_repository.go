package repositories

import (
	"course-management/config"
	"course-management/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(course *models.Course) error
	Update(course *models.Course) error
	Delete(id uuid.UUID) error

	FindByID(id uuid.UUID) (*models.Course, error)
	FindBySlug(slug string) (*models.Course, error)

	FindByInstructor(instructorID uuid.UUID) ([]models.Course, error)
	FindPublished() ([]models.Course, error)

	Search(keyword string, categoryID *uuid.UUID, page int, limit int) ([]models.Course, int64, error)

	CountByInstructor(instructorID uuid.UUID) (int64, error)

	ExistsSlug(slug string) (bool, error)

	UpdateStatus(id uuid.UUID, status string) error
}

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository() CourseRepository {
	return &courseRepository{
		db: config.DB,
	}
}

// Create creates a new course.
func (r *courseRepository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}

// Update updates an existing course.
func (r *courseRepository) Update(course *models.Course) error {
	return r.db.Save(course).Error
}

// Delete soft deletes a course.
func (r *courseRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Course{}, "id = ?", id).Error
}

// FindByID finds a course by ID.
func (r *courseRepository) FindByID(id uuid.UUID) (*models.Course, error) {
	var course models.Course

	err := r.db.
		Preload("Instructor").
		Preload("Category").
		Preload("Sections").
		Preload("Reviews").
		First(&course, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return &course, nil
}

// FindBySlug finds a course by slug.
func (r *courseRepository) FindBySlug(slug string) (*models.Course, error) {
	var course models.Course

	err := r.db.
		Preload("Instructor").
		Preload("Category").
		Preload("Sections").
		Preload("Reviews").
		Where("slug = ?", slug).
		First(&course).
		Error

	if err != nil {
		return nil, err
	}

	return &course, nil
}

// FindByInstructor finds all courses created by an instructor.
func (r *courseRepository) FindByInstructor(
	instructorID uuid.UUID,
) ([]models.Course, error) {

	var courses []models.Course

	err := r.db.
		Preload("Category").
		Where("instructor_id = ?", instructorID).
		Order("created_at DESC").
		Find(&courses).
		Error

	if err != nil {
		return nil, err
	}

	return courses, nil
}

// FindPublished finds all published courses.
func (r *courseRepository) FindPublished() ([]models.Course, error) {

	var courses []models.Course

	err := r.db.
		Preload("Instructor").
		Preload("Category").
		Where("status = ?", "published").
		Order("created_at DESC").
		Find(&courses).
		Error

	if err != nil {
		return nil, err
	}

	return courses, nil
}

// Search searches published courses with optional category filtering
// and pagination.
func (r *courseRepository) Search(
	keyword string,
	categoryID *uuid.UUID,
	page int,
	limit int,
) ([]models.Course, int64, error) {

	var courses []models.Course
	var total int64

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	query := r.db.
		Model(&models.Course{}).
		Where("status = ?", "published")

	if keyword != "" {
		search := "%" + keyword + "%"

		query = query.Where(
			"title ILIKE ? OR description ILIKE ?",
			search,
			search,
		)
	}

	if categoryID != nil {
		query = query.Where(
			"category_id = ?",
			*categoryID,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	err := query.
		Preload("Instructor").
		Preload("Category").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&courses).
		Error

	if err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

// CountByInstructor counts courses belonging to an instructor.
func (r *courseRepository) CountByInstructor(
	instructorID uuid.UUID,
) (int64, error) {

	var count int64

	err := r.db.
		Model(&models.Course{}).
		Where("instructor_id = ?", instructorID).
		Count(&count).
		Error

	return count, err
}

// ExistsSlug checks whether a slug already exists.
func (r *courseRepository) ExistsSlug(slug string) (bool, error) {

	var count int64

	err := r.db.
		Model(&models.Course{}).
		Where("slug = ?", slug).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UpdateStatus updates the status of a course.
func (r *courseRepository) UpdateStatus(
	id uuid.UUID,
	status string,
) error {

	result := r.db.
		Model(&models.Course{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}