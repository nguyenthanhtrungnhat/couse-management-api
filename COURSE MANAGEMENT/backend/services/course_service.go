package services

import (
	"errors"
	"fmt"
	"strings"

	"course-management/dto/course"
	"course-management/dto/lesson"
	"course-management/dto/section"
	"course-management/models"
	"course-management/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseService interface {
	CreateCourse(instructorID uuid.UUID, req course.CreateCourseRequest) (*course.CourseResponse, error)
	GetCourseByID(id uuid.UUID) (*course.CourseResponse, error)
	GetCourseBySlug(slug string) (*course.CourseResponse, error)
	GetMyCourses(instructorID uuid.UUID) ([]course.CourseResponse, error)
	GetPublishedCourses() ([]course.CourseResponse, error)
	SearchCourses(keyword string, categoryID *uuid.UUID, page int, limit int) ([]course.CourseResponse, int64, error)
	UpdateCourse(instructorID uuid.UUID, courseID uuid.UUID, req course.UpdateCourseRequest) (*course.CourseResponse, error)
	DeleteCourse(instructorID uuid.UUID, courseID uuid.UUID) error
}

type courseService struct {
	courseRepository   repositories.CourseRepository
	categoryRepository repositories.CategoryRepository
}

func NewCourseService(
	courseRepository repositories.CourseRepository,
	categoryRepository repositories.CategoryRepository,
) CourseService {
	return &courseService{
		courseRepository:   courseRepository,
		categoryRepository: categoryRepository,
	}
}

// CreateCourse creates a new course for an instructor.
func (s *courseService) CreateCourse(
	instructorID uuid.UUID,
	req course.CreateCourseRequest,
) (*course.CourseResponse, error) {

	title := strings.TrimSpace(req.Title)

	if title == "" {
		return nil, errors.New("course title is required")
	}

	if req.CategoryID == "" {
		return nil, errors.New("category_id is required")
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, errors.New("invalid category_id")
	}

	// Check category.
	_, err = s.categoryRepository.FindByID(categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}

		return nil, err
	}

	// Generate slug.
	slug := generateSlug(title)

	if slug == "" {
		return nil, errors.New("unable to generate course slug")
	}

	// Make sure slug is unique.
	originalSlug := slug
	counter := 1

	for {
		exists, err := s.courseRepository.ExistsSlug(slug)
		if err != nil {
			return nil, err
		}

		if !exists {
			break
		}

		counter++
		slug = fmt.Sprintf("%s-%d", originalSlug, counter)
	}

	description := strings.TrimSpace(req.Description)

	courseModel := &models.Course{
		InstructorID:  instructorID,
		CategoryID:    categoryID,
		Title:         title,
		Slug:          slug,
		Status:        "draft",
		AverageRating: 0,
		TotalStudents: 0,
	}

	// Only set optional description when provided.
	if description != "" {
		courseModel.Description = &description
	}

	if req.ThumbnailURL != nil {
		thumbnail := strings.TrimSpace(*req.ThumbnailURL)

		if thumbnail != "" {
			courseModel.ThumbnailURL = &thumbnail
		}
	}

	if req.PreviewVideoURL != nil {
		previewVideo := strings.TrimSpace(*req.PreviewVideoURL)

		if previewVideo != "" {
			courseModel.PreviewVideoURL = &previewVideo
		}
	}

	courseModel.Price = req.Price

	if err := s.courseRepository.Create(courseModel); err != nil {
		return nil, err
	}

	return s.GetCourseByID(courseModel.ID)
}

// GetCourseByID gets a course by UUID.
func (s *courseService) GetCourseByID(
	id uuid.UUID,
) (*course.CourseResponse, error) {

	courseModel, err := s.courseRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}

		return nil, err
	}

	response := mapCourseToResponse(courseModel)

	return &response, nil
}

// GetCourseBySlug gets a course by slug.
func (s *courseService) GetCourseBySlug(
	slug string,
) (*course.CourseResponse, error) {

	slug = strings.TrimSpace(strings.ToLower(slug))

	if slug == "" {
		return nil, errors.New("slug is required")
	}

	courseModel, err := s.courseRepository.FindBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}

		return nil, err
	}

	response := mapCourseToResponse(courseModel)

	return &response, nil
}

// GetMyCourses gets all courses belonging to an instructor.
func (s *courseService) GetMyCourses(
	instructorID uuid.UUID,
) ([]course.CourseResponse, error) {

	courses, err := s.courseRepository.FindByInstructor(instructorID)
	if err != nil {
		return nil, err
	}

	responses := make([]course.CourseResponse, 0, len(courses))

	for i := range courses {
		response := mapCourseToResponse(&courses[i])
		responses = append(responses, response)
	}

	return responses, nil
}

// GetPublishedCourses gets all published courses.
func (s *courseService) GetPublishedCourses() ([]course.CourseResponse, error) {

	courses, err := s.courseRepository.FindPublished()
	if err != nil {
		return nil, err
	}

	responses := make([]course.CourseResponse, 0, len(courses))

	for i := range courses {
		response := mapCourseToResponse(&courses[i])
		responses = append(responses, response)
	}

	return responses, nil
}

// SearchCourses searches published courses.
func (s *courseService) SearchCourses(
	keyword string,
	categoryID *uuid.UUID,
	page int,
	limit int,
) ([]course.CourseResponse, int64, error) {

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	keyword = strings.TrimSpace(keyword)

	courses, total, err := s.courseRepository.Search(
		keyword,
		categoryID,
		page,
		limit,
	)

	if err != nil {
		return nil, 0, err
	}

	responses := make([]course.CourseResponse, 0, len(courses))

	for i := range courses {
		response := mapCourseToResponse(&courses[i])
		responses = append(responses, response)
	}

	return responses, total, nil
}

// UpdateCourse updates a course owned by an instructor.
func (s *courseService) UpdateCourse(
	instructorID uuid.UUID,
	courseID uuid.UUID,
	req course.UpdateCourseRequest,
) (*course.CourseResponse, error) {

	courseModel, err := s.courseRepository.FindByID(courseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}

		return nil, err
	}

	// Ownership check.
	if courseModel.InstructorID != instructorID {
		return nil, errors.New("you do not have permission to update this course")
	}

	if strings.TrimSpace(req.Title) != "" {
		newTitle := strings.TrimSpace(req.Title)

		// If title changes, regenerate slug.
		if newTitle != courseModel.Title {
			newSlug := generateSlug(newTitle)

			if newSlug == "" {
				return nil, errors.New("unable to generate course slug")
			}

			originalSlug := newSlug
			counter := 1

			for {
				exists, err := s.courseRepository.ExistsSlug(newSlug)
				if err != nil {
					return nil, err
				}

				// Existing slug is allowed if it belongs to this course.
				if exists && newSlug != courseModel.Slug {
					counter++
					newSlug = fmt.Sprintf(
						"%s-%d",
						originalSlug,
						counter,
					)
					continue
				}

				break
			}

			courseModel.Title = newTitle
			courseModel.Slug = newSlug
		}
	}

	if req.CategoryID != "" {
		categoryID, err := uuid.Parse(req.CategoryID)
		if err != nil {
			return nil, errors.New("invalid category_id")
		}

		_, err = s.categoryRepository.FindByID(categoryID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("category not found")
			}

			return nil, err
		}

		courseModel.CategoryID = categoryID
	}

	if req.Description != "" {
		description := strings.TrimSpace(req.Description)
		courseModel.Description = &description
	}

	if req.ThumbnailURL != nil {
		thumbnail := strings.TrimSpace(*req.ThumbnailURL)

		if thumbnail != "" {
			courseModel.ThumbnailURL = &thumbnail
		}
	}

	if req.PreviewVideoURL != nil {
		previewVideo := strings.TrimSpace(*req.PreviewVideoURL)

		if previewVideo != "" {
			courseModel.PreviewVideoURL = &previewVideo
		}
	}

	if req.Price >= 0 {
		courseModel.Price = req.Price
	}

	// Instructor cannot directly publish a course.
	// Keep the current status unchanged here.
	if err := s.courseRepository.Update(courseModel); err != nil {
		return nil, err
	}

	return s.GetCourseByID(courseID)
}

// DeleteCourse deletes a course owned by an instructor.
func (s *courseService) DeleteCourse(
	instructorID uuid.UUID,
	courseID uuid.UUID,
) error {

	courseModel, err := s.courseRepository.FindByID(courseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("course not found")
		}

		return err
	}

	if courseModel.InstructorID != instructorID {
		return errors.New("you do not have permission to delete this course")
	}

	return s.courseRepository.Delete(courseID)
}

// generateSlug creates a URL-friendly slug.
func generateSlug(value string) string {

	value = strings.ToLower(strings.TrimSpace(value))

	var builder strings.Builder

	lastWasDash := false

	for _, char := range value {

		switch {
		case char >= 'a' && char <= 'z':
			builder.WriteRune(char)
			lastWasDash = false

		case char >= '0' && char <= '9':
			builder.WriteRune(char)
			lastWasDash = false

		case char == ' ' || char == '-' || char == '_':
			if !lastWasDash && builder.Len() > 0 {
				builder.WriteRune('-')
				lastWasDash = true
			}
		}
	}

	return strings.Trim(builder.String(), "-")
}

// mapCourseToResponse converts model to DTO.
func mapCourseToResponse(
	c *models.Course,
) course.CourseResponse {

	response := course.CourseResponse{
		ID:            c.ID.String(),
		InstructorID:  c.InstructorID.String(),
		CategoryID:    c.CategoryID.String(),
		Title:         c.Title,
		Slug:          c.Slug,
		Status:        c.Status,
		Price:         c.Price,
		AverageRating: c.AverageRating,
		TotalStudents: c.TotalStudents,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
		Sections:      make([]section.SectionResponse, 0),
	}

	if c.Description != nil {
		response.Description = c.Description
	}

	if c.ThumbnailURL != nil {
		response.ThumbnailURL = c.ThumbnailURL
	}

	if c.PreviewVideoURL != nil {
		response.PreviewVideoURL = c.PreviewVideoURL
	}

	for _, sectionModel := range c.Sections {

		sectionResponse := section.SectionResponse{
			ID:        sectionModel.ID.String(),
			CourseID:  sectionModel.CourseID.String(),
			Title:     sectionModel.Title,
			SortOrder: sectionModel.SortOrder,
			CreatedAt: sectionModel.CreatedAt,
			UpdatedAt: sectionModel.UpdatedAt,
			Lessons:   make([]lesson.LessonResponse, 0),
		}

		for _, lessonModel := range sectionModel.Lessons {

			lessonResponse := lesson.LessonResponse{
				ID:              lessonModel.ID.String(),
				SectionID:       lessonModel.SectionID.String(),
				Title:           lessonModel.Title,
				VideoURL:        lessonModel.VideoURL,
				DurationSeconds: lessonModel.DurationSeconds,
				IsPreview:       lessonModel.IsPreview,
				SortOrder:       lessonModel.SortOrder,
				CreatedAt:       lessonModel.CreatedAt,
				UpdatedAt:       lessonModel.UpdatedAt,
			}

			sectionResponse.Lessons = append(
				sectionResponse.Lessons,
				lessonResponse,
			)
		}

		response.Sections = append(
			response.Sections,
			sectionResponse,
		)
	}

	return response
}
