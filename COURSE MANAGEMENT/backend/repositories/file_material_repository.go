package repositories

import (
	"context"
	"errors"

	"course-management/config"
	"course-management/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type FileMaterialRepository interface {
	Create(material *models.FileMaterial) error
	FindByID(id uuid.UUID) (*models.FileMaterial, error)
	FindByLessonID(lessonID uuid.UUID) ([]models.FileMaterial, error)
	Update(material *models.FileMaterial) error
	Delete(id uuid.UUID) error
}

type fileMaterialRepository struct{}

func NewFileMaterialRepository() FileMaterialRepository {
	return &fileMaterialRepository{}
}

func (r *fileMaterialRepository) Create(
	material *models.FileMaterial,
) error {

	query := `
		INSERT INTO file_materials (
			id,
			lesson_id,
			file_name,
			file_url,
			file_type,
			file_size,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`

	_, err := config.DB.Exec(
		context.Background(),
		query,
		material.ID,
		material.LessonID,
		material.FileName,
		material.FileURL,
		material.FileType,
		material.FileSize,
		material.CreatedAt,
		material.UpdatedAt,
	)

	return err
}

func (r *fileMaterialRepository) FindByID(
	id uuid.UUID,
) (*models.FileMaterial, error) {

	var material models.FileMaterial

	query := `
		SELECT
			id,
			lesson_id,
			file_name,
			file_url,
			file_type,
			file_size,
			created_at,
			updated_at
		FROM file_materials
		WHERE id = $1
	`

	err := config.DB.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&material.ID,
		&material.LessonID,
		&material.FileName,
		&material.FileURL,
		&material.FileType,
		&material.FileSize,
		&material.CreatedAt,
		&material.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}

		return nil, err
	}

	return &material, nil
}

func (r *fileMaterialRepository) FindByLessonID(
	lessonID uuid.UUID,
) ([]models.FileMaterial, error) {

	query := `
		SELECT
			id,
			lesson_id,
			file_name,
			file_url,
			file_type,
			file_size,
			created_at,
			updated_at
		FROM file_materials
		WHERE lesson_id = $1
		ORDER BY created_at ASC
	`

	rows, err := config.DB.Query(
		context.Background(),
		query,
		lessonID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	materials := make([]models.FileMaterial, 0)

	for rows.Next() {

		var material models.FileMaterial

		if err := rows.Scan(
			&material.ID,
			&material.LessonID,
			&material.FileName,
			&material.FileURL,
			&material.FileType,
			&material.FileSize,
			&material.CreatedAt,
			&material.UpdatedAt,
		); err != nil {
			return nil, err
		}

		materials = append(materials, material)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return materials, nil
}

func (r *fileMaterialRepository) Update(
	material *models.FileMaterial,
) error {

	query := `
		UPDATE file_materials
		SET
			file_name = $1,
			file_url = $2,
			file_type = $3,
			file_size = $4,
			updated_at = $5
		WHERE id = $6
	`

	result, err := config.DB.Exec(
		context.Background(),
		query,
		material.FileName,
		material.FileURL,
		material.FileType,
		material.FileSize,
		material.UpdatedAt,
		material.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *fileMaterialRepository) Delete(
	id uuid.UUID,
) error {

	query := `
		DELETE FROM file_materials
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