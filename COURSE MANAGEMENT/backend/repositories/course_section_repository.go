package repositories

import (
	"course-management/config"
	"course-management/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseSectionRepository interface {
	Create(section *models.CourseSection) error
	FindByID(id uuid.UUID) (*models.CourseSection, error)
	FindByCourseID(courseID uuid.UUID) ([]models.CourseSection, error)
	Update(section *models.CourseSection) error
	Delete(id uuid.UUID) error
}

type courseSectionRepository struct {
	db *gorm.DB
}

func NewCourseSectionRepository() CourseSectionRepository {
	return &courseSectionRepository{
		db: config.DB,
	}
}

func (r *courseSectionRepository) Create(
	section *models.CourseSection,
) error {
	return r.db.Create(section).Error
}

func (r *courseSectionRepository) FindByID(
	id uuid.UUID,
) (*models.CourseSection, error) {

	var section models.CourseSection

	err := r.db.
		Preload("Lessons").
		First(&section, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return &section, nil
}

func (r *courseSectionRepository) FindByCourseID(
	courseID uuid.UUID,
) ([]models.CourseSection, error) {

	var sections []models.CourseSection

	err := r.db.
		Where("course_id = ?", courseID).
		Order("sort_order ASC").
		Find(&sections).
		Error

	if err != nil {
		return nil, err
	}

	return sections, nil
}

func (r *courseSectionRepository) Update(
	section *models.CourseSection,
) error {
	return r.db.Save(section).Error
}

func (r *courseSectionRepository) Delete(
	id uuid.UUID,
) error {
	return r.db.Delete(
		&models.CourseSection{},
		"id = ?",
		id,
	).Error
}
