package repositories

import (
	"context"

	"course-management/config"
	"course-management/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type RoleRepository interface {
	FindByID(id string) (*models.Role, error)
	FindByName(name string) (*models.Role, error)
	GetAll() ([]models.Role, error)
}

type roleRepository struct{}

func NewRoleRepository() RoleRepository {
	return &roleRepository{}
}

func (r *roleRepository) FindByID(
	id string,
) (*models.Role, error) {

	roleID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var role models.Role

	query := `
		SELECT
			id,
			name,
			created_at,
			updated_at
		FROM roles
		WHERE id = $1
	`

	err = config.DB.QueryRow(
		context.Background(),
		query,
		roleID,
	).Scan(
		&role.ID,
		&role.Name,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}

		return nil, err
	}

	return &role, nil
}

func (r *roleRepository) FindByName(
	name string,
) (*models.Role, error) {

	var role models.Role

	query := `
		SELECT
			id,
			name,
			created_at,
			updated_at
		FROM roles
		WHERE name = $1
		LIMIT 1
	`

	err := config.DB.QueryRow(
		context.Background(),
		query,
		name,
	).Scan(
		&role.ID,
		&role.Name,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}

		return nil, err
	}

	return &role, nil
}

func (r *roleRepository) GetAll() ([]models.Role, error) {

	query := `
		SELECT
			id,
			name,
			created_at,
			updated_at
		FROM roles
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

	roles := make([]models.Role, 0)

	for rows.Next() {

		var role models.Role

		if err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.CreatedAt,
			&role.UpdatedAt,
		); err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}
