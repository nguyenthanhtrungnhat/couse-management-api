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

	Search(
		keyword string,
		categoryID *uuid.UUID,
		page int,
		limit int,
	) ([]models.Course, int64, error)

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

// ============================================================
// CREATE
// ============================================================

func (r *courseRepository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}

// ============================================================
// UPDATE
// ============================================================

func (r *courseRepository) Update(course *models.Course) error {
	return r.db.Save(course).Error
}

// ============================================================
// DELETE
// ============================================================

func (r *courseRepository) Delete(id uuid.UUID) error {
	return r.db.
		Delete(&models.Course{}, "id = ?", id).
		Error
}

// ============================================================
// PRELOAD COURSE STRUCTURE
// ============================================================

func (r *courseRepository) preloadCourse(query *gorm.DB) *gorm.DB {

	return query.
		Preload("Instructor").
		Preload("Category").
		Preload("Sections").
		Preload("Sections.Lessons").
		Preload("Sections.Lessons.FileMaterials").
		Preload("Reviews")
}

// ============================================================
// FIND BY ID
// ============================================================

func (r *courseRepository) FindByID(
	id uuid.UUID,
) (*models.Course, error) {

	var course models.Course

	err := r.preloadCourse(
		r.db,
	).
		First(&course, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return &course, nil
}

// ============================================================
// FIND BY SLUG
// ============================================================

func (r *courseRepository) FindBySlug(
	slug string,
) (*models.Course, error) {

	var course models.Course

	err := r.preloadCourse(
		r.db,
	).
		Where("slug = ?", slug).
		First(&course).
		Error

	if err != nil {
		return nil, err
	}

	return &course, nil
}

// ============================================================
// FIND BY INSTRUCTOR
// ============================================================

func (r *courseRepository) FindByInstructor(
	instructorID uuid.UUID,
) ([]models.Course, error) {

	var courses []models.Course

	err := r.preloadCourse(
		r.db,
	).
		Where(
			"instructor_id = ?",
			instructorID,
		).
		Order("created_at DESC").
		Find(&courses).
		Error

	if err != nil {
		return nil, err
	}

	return courses, nil
}

// ============================================================
// FIND PUBLISHED
// ============================================================

func (r *courseRepository) FindPublished() ([]models.Course, error) {

	var courses []models.Course

	err := r.preloadCourse(
		r.db,
	).
		Where(
			"status = ?",
			"published",
		).
		Order("created_at DESC").
		Find(&courses).
		Error

	if err != nil {
		return nil, err
	}

	return courses, nil
}

// ============================================================
// SEARCH
// ============================================================

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
		Where(
			"status = ?",
			"published",
		)

	// Keyword search
	if keyword != "" {

		search := "%" + keyword + "%"

		query = query.Where(
			"title ILIKE ? OR description ILIKE ?",
			search,
			search,
		)
	}

	// Category filter
	if categoryID != nil {

		query = query.Where(
			"category_id = ?",
			*categoryID,
		)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	err := r.preloadCourse(
		query,
	).
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

// ============================================================
// COUNT BY INSTRUCTOR
// ============================================================

func (r *courseRepository) CountByInstructor(
	instructorID uuid.UUID,
) (int64, error) {

	var count int64

	err := r.db.
		Model(&models.Course{}).
		Where(
			"instructor_id = ?",
			instructorID,
		).
		Count(&count).
		Error

	return count, err
}

// ============================================================
// EXISTS SLUG
// ============================================================

func (r *courseRepository) ExistsSlug(
	slug string,
) (bool, error) {

	var count int64

	err := r.db.
		Model(&models.Course{}).
		Where(
			"slug = ?",
			slug,
		).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ============================================================
// UPDATE STATUS
// ============================================================

func (r *courseRepository) UpdateStatus(
	id uuid.UUID,
	status string,
) error {

	result := r.db.
		Model(&models.Course{}).
		Where(
			"id = ?",
			id,
		).
		Update(
			"status",
			status,
		)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
