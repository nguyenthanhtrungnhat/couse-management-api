package repositories

import (
	"course-management/config"
	"course-management/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileMaterialRepository interface {
	Create(material *models.FileMaterial) error
	FindByID(id uuid.UUID) (*models.FileMaterial, error)
	FindByLessonID(lessonID uuid.UUID) ([]models.FileMaterial, error)
	Update(material *models.FileMaterial) error
	Delete(id uuid.UUID) error
}

type fileMaterialRepository struct {
	db *gorm.DB
}

func NewFileMaterialRepository() FileMaterialRepository {
	return &fileMaterialRepository{
		db: config.DB,
	}
}

func (r *fileMaterialRepository) Create(
	material *models.FileMaterial,
) error {
	return r.db.Create(material).Error
}

func (r *fileMaterialRepository) FindByID(
	id uuid.UUID,
) (*models.FileMaterial, error) {

	var material models.FileMaterial

	err := r.db.
		First(&material, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return &material, nil
}

func (r *fileMaterialRepository) FindByLessonID(
	lessonID uuid.UUID,
) ([]models.FileMaterial, error) {

	var materials []models.FileMaterial

	err := r.db.
		Where("lesson_id = ?", lessonID).
		Order("created_at ASC").
		Find(&materials).
		Error

	if err != nil {
		return nil, err
	}

	return materials, nil
}

func (r *fileMaterialRepository) Update(
	material *models.FileMaterial,
) error {
	return r.db.Save(material).Error
}

func (r *fileMaterialRepository) Delete(
	id uuid.UUID,
) error {
	return r.db.Delete(
		&models.FileMaterial{},
		"id = ?",
		id,
	).Error
}
