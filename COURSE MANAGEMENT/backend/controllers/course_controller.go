package controllers

import (
	"strconv"

	"course-management/dto/course"
	"course-management/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CourseController struct {
	courseService services.CourseService
}

func NewCourseController(
	courseService services.CourseService,
) *CourseController {
	return &CourseController{
		courseService: courseService,
	}
}

// CreateCourse
// POST /api/courses
func (c *CourseController) CreateCourse(ctx *fiber.Ctx) error {
	var req course.CreateCourseRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid request body",
		})
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	result, err := c.courseService.CreateCourse(userID, req)
	if err != nil {
		return handleCourseError(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "course created successfully",
		"data":    result,
	})
}

// GetCourseByID
// GET /api/courses/:id
func (c *CourseController) GetCourseByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid course id",
		})
	}

	result, err := c.courseService.GetCourseByID(id)
	if err != nil {
		return handleCourseError(ctx, err)
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// GetCourseBySlug
// GET /api/courses/slug/:slug
func (c *CourseController) GetCourseBySlug(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")

	result, err := c.courseService.GetCourseBySlug(slug)
	if err != nil {
		return handleCourseError(ctx, err)
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// GetMyCourses
// GET /api/courses/my
func (c *CourseController) GetMyCourses(ctx *fiber.Ctx) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	result, err := c.courseService.GetMyCourses(userID)
	if err != nil {
		return handleCourseError(ctx, err)
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// GetPublishedCourses
// GET /api/courses/published
func (c *CourseController) GetPublishedCourses(ctx *fiber.Ctx) error {
	result, err := c.courseService.GetPublishedCourses()
	if err != nil {
		return handleCourseError(ctx, err)
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// SearchCourses
// GET /api/courses?keyword=go&category_id=xxx&page=1&limit=10
func (c *CourseController) SearchCourses(ctx *fiber.Ctx) error {
	keyword := ctx.Query("keyword")

	var categoryID *uuid.UUID

	categoryIDString := ctx.Query("category_id")

	if categoryIDString != "" {
		id, err := uuid.Parse(categoryIDString)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "invalid category_id",
			})
		}

		categoryID = &id
	}

	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	result, total, err := c.courseService.SearchCourses(
		keyword,
		categoryID,
		page,
		limit,
	)

	if err != nil {
		return handleCourseError(ctx, err)
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    result,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// UpdateCourse
// PUT /api/courses/:id
func (c *CourseController) UpdateCourse(ctx *fiber.Ctx) error {
	courseID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid course id",
		})
	}

	var req course.UpdateCourseRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid request body",
		})
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	result, err := c.courseService.UpdateCourse(
		userID,
		courseID,
		req,
	)

	if err != nil {
		return handleCourseError(ctx, err)
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "course updated successfully",
		"data":    result,
	})
}

// DeleteCourse
// DELETE /api/courses/:id
func (c *CourseController) DeleteCourse(ctx *fiber.Ctx) error {
	courseID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid course id",
		})
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	if err := c.courseService.DeleteCourse(
		userID,
		courseID,
	); err != nil {
		return handleCourseError(ctx, err)
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "course deleted successfully",
	})
}

// getUserIDFromContext gets authenticated user's UUID.
func getUserIDFromContext(ctx *fiber.Ctx) (uuid.UUID, error) {
	userIDValue := ctx.Locals("userID")

	if userIDValue == nil {
		return uuid.Nil, fiber.NewError(
			fiber.StatusUnauthorized,
			"unauthorized",
		)
	}

	switch value := userIDValue.(type) {
	case uuid.UUID:
		return value, nil

	case string:
		id, err := uuid.Parse(value)
		if err != nil {
			return uuid.Nil, fiber.NewError(
				fiber.StatusUnauthorized,
				"invalid user id",
			)
		}

		return id, nil

	default:
		return uuid.Nil, fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid user id",
		)
	}
}

// handleCourseError converts service errors into HTTP responses.
func handleCourseError(
	ctx *fiber.Ctx,
	err error,
) error {

	message := err.Error()

	switch message {

	case "course not found",
		"category not found":
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": message,
		})

	case "you do not have permission to update this course",
		"you do not have permission to delete this course":
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": message,
		})

	case "unauthorized":
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": message,
		})

	default:
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
		})
	}
}