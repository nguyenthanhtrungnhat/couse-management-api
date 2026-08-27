package services

import (
	"course-management/config"
	"course-management/models"

	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	db *gorm.DB
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

	// --------------------------------------------------------
	// 1. Find lesson -> section -> course
	// --------------------------------------------------------

	var lesson models.Lesson

	err := s.db.
		Preload("Section").
		Preload("Section.Course").
		First(&lesson, "id = ?", lessonID).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lesson not found")
		}

		return nil, err
	}

	// --------------------------------------------------------
	// 2. Check instructor owns the course
	// --------------------------------------------------------

	if lesson.Section.Course.InstructorID != instructorID {
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
		return nil, errors.New(
			"unsupported file type",
		)
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
	// 6. Save file
	// --------------------------------------------------------

	if err := saveUploadedFile(file, filePath); err != nil {
		return nil, err
	}

	// --------------------------------------------------------
	// 7. Get next sort order
	// --------------------------------------------------------

	// --------------------------------------------------------
	// 8. Create database record
	// --------------------------------------------------------

	material := &models.FileMaterial{
		BaseModel: models.BaseModel{
			ID: fileID,
		},

		LessonID: lessonID,

		FileName: file.Filename,

		FileURL: fmt.Sprintf(
			"/uploads/materials/%s",
			filename,
		),

		FileSize: file.Size,

		FileType: ext,
	}

	if err := s.db.Create(material).Error; err != nil {

		// Remove physical file if DB insert fails
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

	var materials []models.FileMaterial

	err := s.db.
		Where("lesson_id = ?", lessonID).
		Order("sort_order ASC").
		Order("created_at ASC").
		Find(&materials).
		Error

	if err != nil {
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

	var material models.FileMaterial

	err := s.db.
		First(&material, "id = ?", materialID).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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

	// --------------------------------------------------------
	// 1. Find material -> lesson -> section -> course
	// --------------------------------------------------------

	var material models.FileMaterial

	err := s.db.
		Preload("Lesson").
		Preload("Lesson.Section").
		Preload("Lesson.Section.Course").
		First(&material, "id = ?", materialID).
		Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("material not found")
		}

		return err
	}

	// --------------------------------------------------------
	// 2. Check instructor ownership
	// --------------------------------------------------------

	if material.Lesson.Section.Course.InstructorID != instructorID {
		return errors.New(
			"you do not have permission to delete this material",
		)
	}

	// --------------------------------------------------------
	// 3. Delete physical file
	// --------------------------------------------------------

	if material.FileURL != "" {

		filePath := strings.TrimPrefix(
			material.FileURL,
			"/",
		)

		_ = os.Remove(filePath)
	}

	// --------------------------------------------------------
	// 4. Delete DB record
	// --------------------------------------------------------

	if err := s.db.Delete(
		&models.FileMaterial{},
		"id = ?",
		materialID,
	).Error; err != nil {

		return err
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

	buffer := make([]byte, 32*1024)

	for {

		n, readErr := src.Read(buffer)

		if n > 0 {

			if _, err := dst.Write(buffer[:n]); err != nil {
				return err
			}
		}

		if readErr != nil {

			if readErr.Error() == "EOF" {
				break
			}

			break
		}
	}

	return nil
}
