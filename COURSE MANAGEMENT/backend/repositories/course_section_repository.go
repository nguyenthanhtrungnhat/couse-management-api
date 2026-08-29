package repositories

import (
	"context"
	"errors"

	"course-management/config"
	"course-management/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CourseSectionRepository interface {
	Create(section *models.CourseSection) error
	FindByID(id uuid.UUID) (*models.CourseSection, error)
	FindByCourseID(courseID uuid.UUID) ([]models.CourseSection, error)
	Update(section *models.CourseSection) error
	Delete(id uuid.UUID) error
}

type courseSectionRepository struct{}

func NewCourseSectionRepository() CourseSectionRepository {
	return &courseSectionRepository{}
}

func (r *courseSectionRepository) Create(
	section *models.CourseSection,
) error {

	query := `
		INSERT INTO course_sections (
			id,
			course_id,
			title,
			sort_order,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := config.DB.Exec(
		context.Background(),
		query,
		section.ID,
		section.CourseID,
		section.Title,
		section.SortOrder,
		section.CreatedAt,
		section.UpdatedAt,
	)

	return err
}

func (r *courseSectionRepository) FindByID(
	id uuid.UUID,
) (*models.CourseSection, error) {

	var section models.CourseSection

	query := `
		SELECT
			id,
			course_id,
			title,
			sort_order,
			created_at,
			updated_at
		FROM course_sections
		WHERE id = $1
	`

	err := config.DB.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&section.ID,
		&section.CourseID,
		&section.Title,
		&section.SortOrder,
		&section.CreatedAt,
		&section.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}

		return nil, err
	}

	return &section, nil
}

func (r *courseSectionRepository) FindByCourseID(
	courseID uuid.UUID,
) ([]models.CourseSection, error) {

	query := `
		SELECT
			id,
			course_id,
			title,
			sort_order,
			created_at,
			updated_at
		FROM course_sections
		WHERE course_id = $1
		ORDER BY sort_order ASC
	`

	rows, err := config.DB.Query(
		context.Background(),
		query,
		courseID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sections := make([]models.CourseSection, 0)

	for rows.Next() {

		var section models.CourseSection

		if err := rows.Scan(
			&section.ID,
			&section.CourseID,
			&section.Title,
			&section.SortOrder,
			&section.CreatedAt,
			&section.UpdatedAt,
		); err != nil {
			return nil, err
		}

		sections = append(sections, section)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sections, nil
}

func (r *courseSectionRepository) Update(
	section *models.CourseSection,
) error {

	query := `
		UPDATE course_sections
		SET
			title = $1,
			sort_order = $2,
			updated_at = $3
		WHERE id = $4
	`

	result, err := config.DB.Exec(
		context.Background(),
		query,
		section.Title,
		section.SortOrder,
		section.UpdatedAt,
		section.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *courseSectionRepository) Delete(
	id uuid.UUID,
) error {

	query := `
		DELETE FROM course_sections
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