package repositories

import (
	"context"
	"errors"

	"course-management/config"
	"course-management/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type LessonRepository interface {
	Create(lesson *models.Lesson) error
	FindByID(id uuid.UUID) (*models.Lesson, error)
	FindBySectionID(sectionID uuid.UUID) ([]models.Lesson, error)
	Update(lesson *models.Lesson) error
	Delete(id uuid.UUID) error
}

type lessonRepository struct{}

func NewLessonRepository() LessonRepository {
	return &lessonRepository{}
}

func (r *lessonRepository) Create(
	lesson *models.Lesson,
) error {

	query := `
		INSERT INTO lessons (
			id,
			section_id,
			title,
			video_url,
			duration_seconds,
			is_preview,
			sort_order,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`

	_, err := config.DB.Exec(
		context.Background(),
		query,
		lesson.ID,
		lesson.SectionID,
		lesson.Title,
		lesson.VideoURL,
		lesson.DurationSeconds,
		lesson.IsPreview,
		lesson.SortOrder,
		lesson.CreatedAt,
		lesson.UpdatedAt,
	)

	return err
}

func (r *lessonRepository) FindByID(
	id uuid.UUID,
) (*models.Lesson, error) {

	var lesson models.Lesson

	query := `
		SELECT
			id,
			section_id,
			title,
			video_url,
			duration_seconds,
			is_preview,
			sort_order,
			created_at,
			updated_at
		FROM lessons
		WHERE id = $1
	`

	err := config.DB.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&lesson.ID,
		&lesson.SectionID,
		&lesson.Title,
		&lesson.VideoURL,
		&lesson.DurationSeconds,
		&lesson.IsPreview,
		&lesson.SortOrder,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}

		return nil, err
	}

	return &lesson, nil
}

func (r *lessonRepository) FindBySectionID(
	sectionID uuid.UUID,
) ([]models.Lesson, error) {

	query := `
		SELECT
			id,
			section_id,
			title,
			video_url,
			duration_seconds,
			is_preview,
			sort_order,
			created_at,
			updated_at
		FROM lessons
		WHERE section_id = $1
		ORDER BY sort_order ASC
	`

	rows, err := config.DB.Query(
		context.Background(),
		query,
		sectionID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	lessons := make([]models.Lesson, 0)

	for rows.Next() {

		var lesson models.Lesson

		if err := rows.Scan(
			&lesson.ID,
			&lesson.SectionID,
			&lesson.Title,
			&lesson.VideoURL,
			&lesson.DurationSeconds,
			&lesson.IsPreview,
			&lesson.SortOrder,
			&lesson.CreatedAt,
			&lesson.UpdatedAt,
		); err != nil {
			return nil, err
		}

		lessons = append(lessons, lesson)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lessons, nil
}

func (r *lessonRepository) Update(
	lesson *models.Lesson,
) error {

	query := `
		UPDATE lessons
		SET
			title = $1,
			video_url = $2,
			duration_seconds = $3,
			is_preview = $4,
			sort_order = $5,
			updated_at = $6
		WHERE id = $7
	`

	result, err := config.DB.Exec(
		context.Background(),
		query,
		lesson.Title,
		lesson.VideoURL,
		lesson.DurationSeconds,
		lesson.IsPreview,
		lesson.SortOrder,
		lesson.UpdatedAt,
		lesson.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *lessonRepository) Delete(
	id uuid.UUID,
) error {

	query := `
		DELETE FROM lessons
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