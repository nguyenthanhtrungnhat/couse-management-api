package repositories

import (
	"course-management/config"
	"course-management/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *models.Category) error
	Update(category *models.Category) error
	Delete(id uuid.UUID) error

	FindByID(id uuid.UUID) (*models.Category, error)
	FindByName(name string) (*models.Category, error)
	FindAll() ([]models.Category, error)

	ExistsName(name string) (bool, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository() CategoryRepository {
	return &categoryRepository{
		db: config.DB,
	}
}

// Create creates a new category.
func (r *categoryRepository) Create(
	category *models.Category,
) error {

	return r.db.Create(category).Error
}

// Update updates an existing category.
func (r *categoryRepository) Update(
	category *models.Category,
) error {

	return r.db.Save(category).Error
}

// Delete soft deletes a category.
func (r *categoryRepository) Delete(
	id uuid.UUID,
) error {

	return r.db.
		Delete(&models.Category{}, "id = ?", id).
		Error
}

// FindByID finds a category by ID.
func (r *categoryRepository) FindByID(
	id uuid.UUID,
) (*models.Category, error) {

	var category models.Category

	err := r.db.
		First(&category, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}

// FindByName finds a category by name.
func (r *categoryRepository) FindByName(
	name string,
) (*models.Category, error) {

	var category models.Category

	err := r.db.
		Where("name = ?", name).
		First(&category).
		Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}

// FindAll returns all categories.
func (r *categoryRepository) FindAll() ([]models.Category, error) {

	var categories []models.Category

	err := r.db.
		Order("name ASC").
		Find(&categories).
		Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}

// ExistsName checks whether a category name already exists.
func (r *categoryRepository) ExistsName(
	name string,
) (bool, error) {

	var count int64

	err := r.db.
		Model(&models.Category{}).
		Where("name = ?", name).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}