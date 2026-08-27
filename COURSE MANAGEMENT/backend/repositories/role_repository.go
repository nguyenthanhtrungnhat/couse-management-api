package repositories

import (
	"course-management/config"
	"course-management/models"

	"gorm.io/gorm"
)

type RoleRepository interface {
	FindByID(id string) (*models.Role, error)
	FindByName(name string) (*models.Role, error)
	GetAll() ([]models.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository() RoleRepository {
	return &roleRepository{
		db: config.DB,
	}
}

func (r *roleRepository) FindByID(id string) (*models.Role, error) {
	var role models.Role

	err := r.db.
		Where("id = ?", id).
		First(&role).Error

	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *roleRepository) FindByName(name string) (*models.Role, error) {
	var role models.Role

	err := r.db.
		Where("name = ?", name).
		First(&role).Error

	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *roleRepository) GetAll() ([]models.Role, error) {
	var roles []models.Role

	err := r.db.
		Order("name ASC").
		Find(&roles).Error

	if err != nil {
		return nil, err
	}

	return roles, nil
}