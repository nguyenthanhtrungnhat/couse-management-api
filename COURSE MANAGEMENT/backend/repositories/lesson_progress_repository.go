package repositories

import (
	"context"
	"errors"

	"course-management/config"
	"course-management/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type LessonProgressRepository interface {
	FindByEnrollmentAndLesson(
		enrollmentID uuid.UUID,
		lessonID uuid.UUID,
	) (*models.LessonProgress, error)

	Create(progress *models.LessonProgress) error

	Update(progress *models.LessonProgress) error

	FindByEnrollment(
		enrollmentID uuid.UUID,
	) ([]models.LessonProgress, error)
}

type lessonProgressRepository struct{}

func NewLessonProgressRepository() LessonProgressRepository {
	return &lessonProgressRepository{}
}

func (r *lessonProgressRepository) FindByEnrollmentAndLesson(
	enrollmentID uuid.UUID,
	lessonID uuid.UUID,
) (*models.LessonProgress, error) {

	var progress models.LessonProgress

	query := `
		SELECT
			id,
			enrollment_id,
			lesson_id,
			completed,
			watched_seconds,
			created_at,
			updated_at
		FROM lesson_progresses
		WHERE enrollment_id = $1
		  AND lesson_id = $2
		LIMIT 1
	`

	err := config.DB.QueryRow(
		context.Background(),
		query,
		enrollmentID,
		lessonID,
	).Scan(
		&progress.ID,
		&progress.EnrollmentID,
		&progress.LessonID,
		&progress.Completed,
		&progress.WatchedSeconds,
		&progress.CreatedAt,
		&progress.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}

		return nil, err
	}

	return &progress, nil
}

func (r *lessonProgressRepository) Create(
	progress *models.LessonProgress,
) error {

	query := `
		INSERT INTO lesson_progresses (
			id,
			enrollment_id,
			lesson_id,
			completed,
			watched_seconds,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`

	_, err := config.DB.Exec(
		context.Background(),
		query,
		progress.ID,
		progress.EnrollmentID,
		progress.LessonID,
		progress.Completed,
		progress.WatchedSeconds,
		progress.CreatedAt,
		progress.UpdatedAt,
	)

	return err
}

func (r *lessonProgressRepository) Update(
	progress *models.LessonProgress,
) error {

	query := `
		UPDATE lesson_progresses
		SET
			completed = $1,
			watched_seconds = $2,
			updated_at = $3
		WHERE id = $4
	`

	result, err := config.DB.Exec(
		context.Background(),
		query,
		progress.Completed,
		progress.WatchedSeconds,
		progress.UpdatedAt,
		progress.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *lessonProgressRepository) FindByEnrollment(
	enrollmentID uuid.UUID,
) ([]models.LessonProgress, error) {

	query := `
		SELECT
			id,
			enrollment_id,
			lesson_id,
			completed,
			watched_seconds,
			created_at,
			updated_at
		FROM lesson_progresses
		WHERE enrollment_id = $1
		ORDER BY created_at ASC
	`

	rows, err := config.DB.Query(
		context.Background(),
		query,
		enrollmentID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	progresses := make([]models.LessonProgress, 0)

	for rows.Next() {

		var progress models.LessonProgress

		if err := rows.Scan(
			&progress.ID,
			&progress.EnrollmentID,
			&progress.LessonID,
			&progress.Completed,
			&progress.WatchedSeconds,
			&progress.CreatedAt,
			&progress.UpdatedAt,
		); err != nil {
			return nil, err
		}

		progresses = append(progresses, progress)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return progresses, nil
}
