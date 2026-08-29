package repositories

import (
	"context"
	"errors"

	"course-management/config"
	"course-management/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uuid.UUID) (*models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(
	user *models.User,
) error {

	query := `
		INSERT INTO users (
			id,
			role_id,
			full_name,
			email,
			password_hash,
			avatar_url,
			provider,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`

	_, err := config.DB.Exec(
		context.Background(),
		query,
		user.ID,
		user.RoleID,
		user.FullName,
		user.Email,
		user.PasswordHash,
		user.AvatarURL,
		user.Provider,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (r *userRepository) FindByEmail(
	email string,
) (*models.User, error) {

	var user models.User

	query := `
		SELECT
			u.id,
			u.role_id,
			u.full_name,
			u.email,
			u.password_hash,
			u.avatar_url,
			u.provider,
			u.created_at,
			u.updated_at
		FROM users u
		WHERE u.email = $1
		LIMIT 1
	`

	err := config.DB.QueryRow(
		context.Background(),
		query,
		email,
	).Scan(
		&user.ID,
		&user.RoleID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.AvatarURL,
		&user.Provider,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}

		return nil, err
	}

	// Load Role.
	roleQuery := `
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
		roleQuery,
		user.RoleID,
	).Scan(
		&user.Role.ID,
		&user.Role.Name,
		&user.Role.CreatedAt,
		&user.Role.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByID(
	id uuid.UUID,
) (*models.User, error) {

	var user models.User

	query := `
		SELECT
			id,
			role_id,
			full_name,
			email,
			password_hash,
			avatar_url,
			provider,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	err := config.DB.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&user.ID,
		&user.RoleID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.AvatarURL,
		&user.Provider,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}

		return nil, err
	}

	roleQuery := `
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
		roleQuery,
		user.RoleID,
	).Scan(
		&user.Role.ID,
		&user.Role.Name,
		&user.Role.CreatedAt,
		&user.Role.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Update(
	user *models.User,
) error {

	query := `
		UPDATE users
		SET
			role_id = $1,
			full_name = $2,
			email = $3,
			password_hash = $4,
			avatar_url = $5,
			provider = $6,
			updated_at = $7
		WHERE id = $8
	`

	result, err := config.DB.Exec(
		context.Background(),
		query,
		user.RoleID,
		user.FullName,
		user.Email,
		user.PasswordHash,
		user.AvatarURL,
		user.Provider,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *userRepository) Delete(
	id uuid.UUID,
) error {

	query := `
		DELETE FROM users
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
