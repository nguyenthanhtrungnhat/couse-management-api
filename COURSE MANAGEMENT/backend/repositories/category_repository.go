package repositories

import (
	"context"
	"errors"

	"course-management/config"
	"course-management/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

type categoryRepository struct{}

func NewCategoryRepository() CategoryRepository {
	return &categoryRepository{}
}

func (r *categoryRepository) Create(
	category *models.Category,
) error {

	query := `
		INSERT INTO categories (
			id,
			name,
			description,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := config.DB.Exec(
		context.Background(),
		query,
		category.ID,
		category.Name,
		category.Description,
		category.CreatedAt,
		category.UpdatedAt,
	)

	return err
}

func (r *categoryRepository) Update(
	category *models.Category,
) error {

	query := `
		UPDATE categories
		SET
			name = $1,
			description = $2,
			updated_at = $3
		WHERE id = $4
	`

	result, err := config.DB.Exec(
		context.Background(),
		query,
		category.Name,
		category.Description,
		category.UpdatedAt,
		category.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *categoryRepository) Delete(
	id uuid.UUID,
) error {

	query := `
		DELETE FROM categories
		WHERE id = $1
	`

	result, err := config.DB.Exec(
		context.Background(),
		query,
		id,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *categoryRepository) FindByID(
	id uuid.UUID,
) (*models.Category, error) {

	var category models.Category

	query := `
		SELECT
			id,
			name,
			description,
			created_at,
			updated_at
		FROM categories
		WHERE id = $1
	`

	err := config.DB.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}

		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) FindByName(
	name string,
) (*models.Category, error) {

	var category models.Category

	query := `
		SELECT
			id,
			name,
			description,
			created_at,
			updated_at
		FROM categories
		WHERE name = $1
		LIMIT 1
	`

	err := config.DB.QueryRow(
		context.Background(),
		query,
		name,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}

		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) FindAll() ([]models.Category, error) {

	query := `
		SELECT
			id,
			name,
			description,
			created_at,
			updated_at
		FROM categories
		ORDER BY name ASC
	`

	rows, err := config.DB.Query(
		context.Background(),
		query,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	categories := make([]models.Category, 0)

	for rows.Next() {

		var category models.Category

		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		); err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) ExistsName(
	name string,
) (bool, error) {

	var exists bool

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM categories
			WHERE name = $1
		)
	`

	err := config.DB.QueryRow(
		context.Background(),
		query,
		name,
	).Scan(&exists)

	return exists, err
}