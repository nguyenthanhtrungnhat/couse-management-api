package repositories

import (
	"course-management/config"
	"course-management/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LessonRepository interface {
	Create(lesson *models.Lesson) error
	FindByID(id uuid.UUID) (*models.Lesson, error)
	FindBySectionID(sectionID uuid.UUID) ([]models.Lesson, error)
	Update(lesson *models.Lesson) error
	Delete(id uuid.UUID) error
}

type lessonRepository struct {
	db *gorm.DB
}

func NewLessonRepository() LessonRepository {
	return &lessonRepository{
		db: config.DB,
	}
}

func (r *lessonRepository) Create(
	lesson *models.Lesson,
) error {
	return r.db.Create(lesson).Error
}

func (r *lessonRepository) FindByID(
	id uuid.UUID,
) (*models.Lesson, error) {

	var lesson models.Lesson

	err := r.db.
		Preload("FileMaterials").
		First(&lesson, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return &lesson, nil
}

func (r *lessonRepository) FindBySectionID(
	sectionID uuid.UUID,
) ([]models.Lesson, error) {

	var lessons []models.Lesson

	err := r.db.
		Where("section_id = ?", sectionID).
		Order("sort_order ASC").
		Find(&lessons).
		Error

	if err != nil {
		return nil, err
	}

	return lessons, nil
}

func (r *lessonRepository) Update(
	lesson *models.Lesson,
) error {
	return r.db.Save(lesson).Error
}

func (r *lessonRepository) Delete(
	id uuid.UUID,
) error {
	return r.db.Delete(
		&models.Lesson{},
		"id = ?",
		id,
	).Error
}
