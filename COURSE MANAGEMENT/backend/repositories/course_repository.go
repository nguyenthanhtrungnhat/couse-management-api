package repositories

import (
	"context"
	"errors"

	"course-management/config"
	"course-management/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CourseRepository interface {
	Create(course *models.Course) error
	Update(course *models.Course) error
	Delete(id uuid.UUID) error

	FindByID(id uuid.UUID) (*models.Course, error)
	FindBySlug(slug string) (*models.Course, error)

	FindByInstructor(instructorID uuid.UUID) ([]models.Course, error)
	FindPublished() ([]models.Course, error)

	Search(
		keyword string,
		categoryID *uuid.UUID,
		page int,
		limit int,
	) ([]models.Course, int64, error)

	CountByInstructor(instructorID uuid.UUID) (int64, error)

	ExistsSlug(slug string) (bool, error)

	UpdateStatus(id uuid.UUID, status string) error
}

type courseRepository struct{}

func NewCourseRepository() CourseRepository {
	return &courseRepository{}
}

func (r *courseRepository) Create(
	course *models.Course,
) error {

	query := `
		INSERT INTO courses (
			id,
			instructor_id,
			category_id,
			title,
			slug,
			description,
			thumbnail_url,
			preview_video_url,
			price,
			status,
			average_rating,
			total_students,
			created_at,
			updated_at
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14
		)
	`

	_, err := config.DB.Exec(
		context.Background(),
		query,
		course.ID,
		course.InstructorID,
		course.CategoryID,
		course.Title,
		course.Slug,
		course.Description,
		course.ThumbnailURL,
		course.PreviewVideoURL,
		course.Price,
		course.Status,
		course.AverageRating,
		course.TotalStudents,
		course.CreatedAt,
		course.UpdatedAt,
	)

	return err
}

func (r *courseRepository) Update(
	course *models.Course,
) error {

	query := `
		UPDATE courses
		SET
			instructor_id = $1,
			category_id = $2,
			title = $3,
			slug = $4,
			description = $5,
			thumbnail_url = $6,
			preview_video_url = $7,
			price = $8,
			status = $9,
			average_rating = $10,
			total_students = $11,
			updated_at = $12
		WHERE id = $13
	`

	result, err := config.DB.Exec(
		context.Background(),
		query,
		course.InstructorID,
		course.CategoryID,
		course.Title,
		course.Slug,
		course.Description,
		course.ThumbnailURL,
		course.PreviewVideoURL,
		course.Price,
		course.Status,
		course.AverageRating,
		course.TotalStudents,
		course.UpdatedAt,
		course.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *courseRepository) Delete(
	id uuid.UUID,
) error {

	query := `
		DELETE FROM courses
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

func (r *courseRepository) FindByID(
	id uuid.UUID,
) (*models.Course, error) {

	var course models.Course

	query := `
		SELECT
			id,
			instructor_id,
			category_id,
			title,
			slug,
			description,
			thumbnail_url,
			preview_video_url,
			price,
			status,
			average_rating,
			total_students,
			created_at,
			updated_at
		FROM courses
		WHERE id = $1
	`

	err := config.DB.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&course.ID,
		&course.InstructorID,
		&course.CategoryID,
		&course.Title,
		&course.Slug,
		&course.Description,
		&course.ThumbnailURL,
		&course.PreviewVideoURL,
		&course.Price,
		&course.Status,
		&course.AverageRating,
		&course.TotalStudents,
		&course.CreatedAt,
		&course.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}

		return nil, err
	}

	return r.loadCourseRelations(&course)
}

func (r *courseRepository) FindBySlug(
	slug string,
) (*models.Course, error) {

	var course models.Course

	query := `
		SELECT
			id,
			instructor_id,
			category_id,
			title,
			slug,
			description,
			thumbnail_url,
			preview_video_url,
			price,
			status,
			average_rating,
			total_students,
			created_at,
			updated_at
		FROM courses
		WHERE slug = $1
		LIMIT 1
	`

	err := config.DB.QueryRow(
		context.Background(),
		query,
		slug,
	).Scan(
		&course.ID,
		&course.InstructorID,
		&course.CategoryID,
		&course.Title,
		&course.Slug,
		&course.Description,
		&course.ThumbnailURL,
		&course.PreviewVideoURL,
		&course.Price,
		&course.Status,
		&course.AverageRating,
		&course.TotalStudents,
		&course.CreatedAt,
		&course.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}

		return nil, err
	}

	return r.loadCourseRelations(&course)
}

func (r *courseRepository) FindByInstructor(
	instructorID uuid.UUID,
) ([]models.Course, error) {

	query := `
		SELECT
			id,
			instructor_id,
			category_id,
			title,
			slug,
			description,
			thumbnail_url,
			preview_video_url,
			price,
			status,
			average_rating,
			total_students,
			created_at,
			updated_at
		FROM courses
		WHERE instructor_id = $1
		ORDER BY created_at DESC
	`

	return r.queryCourses(query, instructorID)
}

func (r *courseRepository) FindPublished() ([]models.Course, error) {

	query := `
		SELECT
			id,
			instructor_id,
			category_id,
			title,
			slug,
			description,
			thumbnail_url,
			preview_video_url,
			price,
			status,
			average_rating,
			total_students,
			created_at,
			updated_at
		FROM courses
		WHERE status = $1
		ORDER BY created_at DESC
	`

	return r.queryCourses(query, "published")
}

func (r *courseRepository) Search(
	keyword string,
	categoryID *uuid.UUID,
	page int,
	limit int,
) ([]models.Course, int64, error) {

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	args := []interface{}{"published"}
	argIndex := 2

	query := `
		SELECT
			id,
			instructor_id,
			category_id,
			title,
			slug,
			description,
			thumbnail_url,
			preview_video_url,
			price,
			status,
			average_rating,
			total_students,
			created_at,
			updated_at
		FROM courses
		WHERE status = $1
	`

	countQuery := `
		SELECT COUNT(*)
		FROM courses
		WHERE status = $1
	`

	if keyword != "" {
		search := "%" + keyword + "%"

		query += `
			AND (
				title ILIKE $` + itoa(argIndex) + `
				OR description ILIKE $` + itoa(argIndex) + `
			)
		`

		countQuery += `
			AND (
				title ILIKE $` + itoa(argIndex) + `
				OR description ILIKE $` + itoa(argIndex) + `
			)
		`

		args = append(args, search)
		argIndex++
	}

	if categoryID != nil {
		query += `
			AND category_id = $` + itoa(argIndex) + `
		`

		countQuery += `
			AND category_id = $` + itoa(argIndex) + `
		`

		args = append(args, *categoryID)
		argIndex++
	}

	var total int64

	if err := config.DB.QueryRow(
		context.Background(),
		countQuery,
		args...,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	query += `
		ORDER BY created_at DESC
		LIMIT $` + itoa(argIndex) + `
		OFFSET $` + itoa(argIndex+1)

	args = append(args, limit, offset)

	courses, err := r.queryCourses(
		query,
		args...,
	)

	if err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

func (r *courseRepository) CountByInstructor(
	instructorID uuid.UUID,
) (int64, error) {

	var count int64

	query := `
		SELECT COUNT(*)
		FROM courses
		WHERE instructor_id = $1
	`

	err := config.DB.QueryRow(
		context.Background(),
		query,
		instructorID,
	).Scan(&count)

	return count, err
}

func (r *courseRepository) ExistsSlug(
	slug string,
) (bool, error) {

	var exists bool

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM courses
			WHERE slug = $1
		)
	`

	err := config.DB.QueryRow(
		context.Background(),
		query,
		slug,
	).Scan(&exists)

	return exists, err
}

func (r *courseRepository) UpdateStatus(
	id uuid.UUID,
	status string,
) error {

	query := `
		UPDATE courses
		SET
			status = $1,
			updated_at = NOW()
		WHERE id = $2
	`

	result, err := config.DB.Exec(
		context.Background(),
		query,
		status,
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

func (r *courseRepository) queryCourses(
	query string,
	args ...interface{},
) ([]models.Course, error) {

	rows, err := config.DB.Query(
		context.Background(),
		query,
		args...,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	courses := make([]models.Course, 0)

	for rows.Next() {

		var course models.Course

		if err := rows.Scan(
			&course.ID,
			&course.InstructorID,
			&course.CategoryID,
			&course.Title,
			&course.Slug,
			&course.Description,
			&course.ThumbnailURL,
			&course.PreviewVideoURL,
			&course.Price,
			&course.Status,
			&course.AverageRating,
			&course.TotalStudents,
			&course.CreatedAt,
			&course.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if _, err := r.loadCourseRelations(&course); err != nil {
			return nil, err
		}

		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

func (r *courseRepository) loadCourseRelations(
	course *models.Course,
) (*models.Course, error) {

	// Instructor
	instructorQuery := `
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
		instructorQuery,
		course.InstructorID,
	).Scan(
		&course.Instructor.ID,
		&course.Instructor.RoleID,
		&course.Instructor.FullName,
		&course.Instructor.Email,
		&course.Instructor.PasswordHash,
		&course.Instructor.AvatarURL,
		&course.Instructor.Provider,
		&course.Instructor.CreatedAt,
		&course.Instructor.UpdatedAt,
	)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// Category
	categoryQuery := `
		SELECT
			id,
			name,
			description,
			created_at,
			updated_at
		FROM categories
		WHERE id = $1
	`

	err = config.DB.QueryRow(
		context.Background(),
		categoryQuery,
		course.CategoryID,
	).Scan(
		&course.Category.ID,
		&course.Category.Name,
		&course.Category.Description,
		&course.Category.CreatedAt,
		&course.Category.UpdatedAt,
	)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// Sections
	sectionQuery := `
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

	sectionRows, err := config.DB.Query(
		context.Background(),
		sectionQuery,
		course.ID,
	)

	if err != nil {
		return nil, err
	}

	defer sectionRows.Close()

	course.Sections = make([]models.CourseSection, 0)

	for sectionRows.Next() {

		var section models.CourseSection

		if err := sectionRows.Scan(
			&section.ID,
			&section.CourseID,
			&section.Title,
			&section.SortOrder,
			&section.CreatedAt,
			&section.UpdatedAt,
		); err != nil {
			return nil, err
		}

		lessonQuery := `
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

		lessonRows, err := config.DB.Query(
			context.Background(),
			lessonQuery,
			section.ID,
		)

		if err != nil {
			return nil, err
		}

		section.Lessons = make([]models.Lesson, 0)

		for lessonRows.Next() {

			var lesson models.Lesson

			if err := lessonRows.Scan(
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
				lessonRows.Close()
				return nil, err
			}

			materialQuery := `
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

			materialRows, err := config.DB.Query(
				context.Background(),
				materialQuery,
				lesson.ID,
			)

			if err != nil {
				lessonRows.Close()
				return nil, err
			}

			lesson.FileMaterials = make(
				[]models.FileMaterial,
				0,
			)

			for materialRows.Next() {

				var material models.FileMaterial

				if err := materialRows.Scan(
					&material.ID,
					&material.LessonID,
					&material.FileName,
					&material.FileURL,
					&material.FileType,
					&material.FileSize,
					&material.CreatedAt,
					&material.UpdatedAt,
				); err != nil {
					materialRows.Close()
					lessonRows.Close()
					return nil, err
				}

				lesson.FileMaterials = append(
					lesson.FileMaterials,
					material,
				)
			}

			materialRows.Close()

			if err := materialRows.Err(); err != nil {
				lessonRows.Close()
				return nil, err
			}

			section.Lessons = append(
				section.Lessons,
				lesson,
			)
		}

		lessonRows.Close()

		if err := lessonRows.Err(); err != nil {
			return nil, err
		}

		course.Sections = append(
			course.Sections,
			section,
		)
	}

	if err := sectionRows.Err(); err != nil {
		return nil, err
	}

	return course, nil
}

func itoa(value int) string {
	switch value {
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	case 4:
		return "4"
	case 5:
		return "5"
	case 6:
		return "6"
	case 7:
		return "7"
	case 8:
		return "8"
	case 9:
		return "9"
	case 10:
		return "10"
	case 11:
		return "11"
	case 12:
		return "12"
	case 13:
		return "13"
	case 14:
		return "14"
	case 15:
		return "15"
	case 16:
		return "16"
	default:
		return "17"
	}
}
