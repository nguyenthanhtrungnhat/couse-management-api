package repositories

import (
	"context"
	"errors"

	"course-management/config"
	"course-management/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type EnrollmentRepository interface {
	Create(enrollment *models.Enrollment) error
	FindByID(id uuid.UUID) (*models.Enrollment, error)
	FindByUserAndCourse(userID uuid.UUID, courseID uuid.UUID) (*models.Enrollment, error)
	FindByUser(userID uuid.UUID) ([]models.Enrollment, error)
	FindByCourse(courseID uuid.UUID) ([]models.Enrollment, error)
	Exists(userID uuid.UUID, courseID uuid.UUID) (bool, error)
	Delete(id uuid.UUID) error
	CountByCourse(courseID uuid.UUID) (int64, error)
}

type enrollmentRepository struct{}

func NewEnrollmentRepository() EnrollmentRepository {
	return &enrollmentRepository{}
}

func (r *enrollmentRepository) Create(
	enrollment *models.Enrollment,
) error {

	query := `
		INSERT INTO enrollments (
			id,
			user_id,
			course_id,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4,$5)
	`

	_, err := config.DB.Exec(
		context.Background(),
		query,
		enrollment.ID,
		enrollment.UserID,
		enrollment.CourseID,
		enrollment.CreatedAt,
		enrollment.UpdatedAt,
	)

	return err
}

func (r *enrollmentRepository) FindByID(
	id uuid.UUID,
) (*models.Enrollment, error) {

	var enrollment models.Enrollment

	query := `
		SELECT
			id,
			user_id,
			course_id,
			created_at,
			updated_at
		FROM enrollments
		WHERE id = $1
	`

	err := config.DB.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&enrollment.ID,
		&enrollment.UserID,
		&enrollment.CourseID,
		&enrollment.CreatedAt,
		&enrollment.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}

		return nil, err
	}

	return &enrollment, nil
}

func (r *enrollmentRepository) FindByUserAndCourse(
	userID uuid.UUID,
	courseID uuid.UUID,
) (*models.Enrollment, error) {

	var enrollment models.Enrollment

	query := `
		SELECT
			id,
			user_id,
			course_id,
			created_at,
			updated_at
		FROM enrollments
		WHERE user_id = $1
		  AND course_id = $2
		LIMIT 1
	`

	err := config.DB.QueryRow(
		context.Background(),
		query,
		userID,
		courseID,
	).Scan(
		&enrollment.ID,
		&enrollment.UserID,
		&enrollment.CourseID,
		&enrollment.CreatedAt,
		&enrollment.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}

		return nil, err
	}

	return &enrollment, nil
}

func (r *enrollmentRepository) FindByUser(
	userID uuid.UUID,
) ([]models.Enrollment, error) {

	query := `
		SELECT
			id,
			user_id,
			course_id,
			created_at,
			updated_at
		FROM enrollments
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := config.DB.Query(
		context.Background(),
		query,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	enrollments := make([]models.Enrollment, 0)

	for rows.Next() {

		var enrollment models.Enrollment

		if err := rows.Scan(
			&enrollment.ID,
			&enrollment.UserID,
			&enrollment.CourseID,
			&enrollment.CreatedAt,
			&enrollment.UpdatedAt,
		); err != nil {
			return nil, err
		}

		enrollments = append(enrollments, enrollment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return enrollments, nil
}

func (r *enrollmentRepository) FindByCourse(
	courseID uuid.UUID,
) ([]models.Enrollment, error) {

	query := `
		SELECT
			id,
			user_id,
			course_id,
			created_at,
			updated_at
		FROM enrollments
		WHERE course_id = $1
		ORDER BY created_at DESC
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

	enrollments := make([]models.Enrollment, 0)

	for rows.Next() {

		var enrollment models.Enrollment

		if err := rows.Scan(
			&enrollment.ID,
			&enrollment.UserID,
			&enrollment.CourseID,
			&enrollment.CreatedAt,
			&enrollment.UpdatedAt,
		); err != nil {
			return nil, err
		}

		enrollments = append(enrollments, enrollment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return enrollments, nil
}

func (r *enrollmentRepository) Exists(
	userID uuid.UUID,
	courseID uuid.UUID,
) (bool, error) {

	var exists bool

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM enrollments
			WHERE user_id = $1
			  AND course_id = $2
		)
	`

	err := config.DB.QueryRow(
		context.Background(),
		query,
		userID,
		courseID,
	).Scan(&exists)

	return exists, err
}

func (r *enrollmentRepository) Delete(
	id uuid.UUID,
) error {

	query := `
		DELETE FROM enrollments
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

func (r *enrollmentRepository) CountByCourse(
	courseID uuid.UUID,
) (int64, error) {

	var count int64

	query := `
		SELECT COUNT(*)
		FROM enrollments
		WHERE course_id = $1
	`

	err := config.DB.QueryRow(
		context.Background(),
		query,
		courseID,
	).Scan(&count)

	return count, err
}