package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"course-management/config"
	"course-management/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FileMaterialService interface {
	CreateMaterial(
		instructorID uuid.UUID,
		lessonID uuid.UUID,
		file *multipart.FileHeader,
	) (*models.FileMaterial, error)

	GetMaterialsByLesson(
		lessonID uuid.UUID,
	) ([]models.FileMaterial, error)

	GetMaterialByID(
		materialID uuid.UUID,
	) (*models.FileMaterial, error)

	DeleteMaterial(
		instructorID uuid.UUID,
		materialID uuid.UUID,
	) error
}

type fileMaterialService struct {
	db *pgxpool.Pool
}

func NewFileMaterialService() FileMaterialService {
	return &fileMaterialService{
		db: config.DB,
	}
}

// ============================================================
// CREATE MATERIAL
// ============================================================

func (s *fileMaterialService) CreateMaterial(
	instructorID uuid.UUID,
	lessonID uuid.UUID,
	file *multipart.FileHeader,
) (*models.FileMaterial, error) {

	if file == nil {
		return nil, errors.New("file is required")
	}

	if file.Size <= 0 {
		return nil, errors.New("file is empty")
	}

	ctx := context.Background()

	// --------------------------------------------------------
	// 1. Find lesson -> section -> course
	// --------------------------------------------------------

	var courseInstructorID uuid.UUID

	err := s.db.QueryRow(
		ctx,
		`SELECT c.instructor_id
		 FROM lessons l
		 JOIN course_sections cs
		   ON cs.id = l.section_id
		 JOIN courses c
		   ON c.id = cs.course_id
		 WHERE l.id = $1`,
		lessonID,
	).Scan(&courseInstructorID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("lesson not found")
		}

		return nil, err
	}

	// --------------------------------------------------------
	// 2. Check instructor owns the course
	// --------------------------------------------------------

	if courseInstructorID != instructorID {
		return nil, errors.New(
			"you do not have permission to add material to this lesson",
		)
	}

	// --------------------------------------------------------
	// 3. Validate file extension
	// --------------------------------------------------------

	ext := strings.ToLower(filepath.Ext(file.Filename))

	allowedExtensions := map[string]bool{
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".ppt":  true,
		".pptx": true,
		".xls":  true,
		".xlsx": true,
		".txt":  true,
		".zip":  true,
		".rar":  true,
	}

	if !allowedExtensions[ext] {
		return nil, errors.New("unsupported file type")
	}

	// --------------------------------------------------------
	// 4. Create upload directory
	// --------------------------------------------------------

	uploadDir := "./uploads/materials"

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, err
	}

	// --------------------------------------------------------
	// 5. Generate unique filename
	// --------------------------------------------------------

	fileID := uuid.New()

	filename := fmt.Sprintf(
		"%s%s",
		fileID.String(),
		ext,
	)

	filePath := filepath.Join(
		uploadDir,
		filename,
	)

	// --------------------------------------------------------
	// 6. Save physical file
	// --------------------------------------------------------

	if err := saveUploadedFile(file, filePath); err != nil {
		return nil, err
	}

	// --------------------------------------------------------
	// 7. Create database record
	// --------------------------------------------------------

	fileURL := fmt.Sprintf(
		"/uploads/materials/%s",
		filename,
	)

	fileType := ext

	material := &models.FileMaterial{
		BaseModel: models.BaseModel{
			ID: fileID,
		},

		LessonID: lessonID,
		FileName: file.Filename,
		FileURL:  &fileURL,
		FileSize: file.Size,
		FileType: &fileType,
	}

	_, err = s.db.Exec(
		ctx,
		`INSERT INTO file_materials (
			id,
			lesson_id,
			file_name,
			file_url,
			file_type,
			file_size,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`,
		material.ID,
		material.LessonID,
		material.FileName,
		material.FileURL,
		material.FileType,
		material.FileSize,
	)

	if err != nil {
		// Remove physical file if DB insert fails.
		_ = os.Remove(filePath)

		return nil, err
	}

	return material, nil
}

// ============================================================
// GET MATERIALS BY LESSON
// ============================================================

func (s *fileMaterialService) GetMaterialsByLesson(
	lessonID uuid.UUID,
) ([]models.FileMaterial, error) {

	ctx := context.Background()

	rows, err := s.db.Query(
		ctx,
		`SELECT
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
		ORDER BY created_at ASC`,
		lessonID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	materials := make([]models.FileMaterial, 0)

	for rows.Next() {
		var material models.FileMaterial

		err := rows.Scan(
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
			return nil, err
		}

		materials = append(materials, material)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return materials, nil
}

// ============================================================
// GET MATERIAL BY ID
// ============================================================

func (s *fileMaterialService) GetMaterialByID(
	materialID uuid.UUID,
) (*models.FileMaterial, error) {

	ctx := context.Background()

	var material models.FileMaterial

	err := s.db.QueryRow(
		ctx,
		`SELECT
			id,
			lesson_id,
			file_name,
			file_url,
			file_type,
			file_size,
			created_at,
			updated_at
		FROM file_materials
		WHERE id = $1`,
		materialID,
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
			return nil, errors.New("material not found")
		}

		return nil, err
	}

	return &material, nil
}

// ============================================================
// DELETE MATERIAL
// ============================================================

func (s *fileMaterialService) DeleteMaterial(
	instructorID uuid.UUID,
	materialID uuid.UUID,
) error {

	ctx := context.Background()

	// --------------------------------------------------------
	// 1. Find material -> lesson -> section -> course
	// --------------------------------------------------------

	var (
		lessonID           uuid.UUID
		fileURL            string
		courseInstructorID uuid.UUID
	)

	err := s.db.QueryRow(
		ctx,
		`SELECT
			fm.lesson_id,
			fm.file_url,
			c.instructor_id
		FROM file_materials fm
		JOIN lessons l
		  ON l.id = fm.lesson_id
		JOIN course_sections cs
		  ON cs.id = l.section_id
		JOIN courses c
		  ON c.id = cs.course_id
		WHERE fm.id = $1`,
		materialID,
	).Scan(
		&lessonID,
		&fileURL,
		&courseInstructorID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("material not found")
		}

		return err
	}

	_ = lessonID

	// --------------------------------------------------------
	// 2. Check instructor ownership
	// --------------------------------------------------------

	if courseInstructorID != instructorID {
		return errors.New(
			"you do not have permission to delete this material",
		)
	}

	// --------------------------------------------------------
	// 3. Delete database record
	// --------------------------------------------------------

	result, err := s.db.Exec(
		ctx,
		`DELETE FROM file_materials
		 WHERE id = $1`,
		materialID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("material not found")
	}

	// --------------------------------------------------------
	// 4. Delete physical file
	// --------------------------------------------------------

	if fileURL != "" {
		filePath := strings.TrimPrefix(fileURL, "/")

		_ = os.Remove(filePath)
	}

	return nil
}

// ============================================================
// SAVE UPLOADED FILE
// ============================================================

func saveUploadedFile(
	file *multipart.FileHeader,
	destination string,
) error {

	src, err := file.Open()
	if err != nil {
		return err
	}

	defer src.Close()

	dst, err := os.Create(destination)
	if err != nil {
		return err
	}

	defer dst.Close()

	_, err = io.Copy(dst, src)

	return err
}
