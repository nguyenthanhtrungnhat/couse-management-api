package controllers

import (
	"errors"

	"course-management/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type EnrollmentController struct {
	service services.EnrollmentService
}

func NewEnrollmentController(
	service services.EnrollmentService,
) *EnrollmentController {
	return &EnrollmentController{
		service: service,
	}
}

func getUserID(ctx *fiber.Ctx) (uuid.UUID, error) {

	value := ctx.Locals("user_id")

	userID, ok := value.(string)

	if !ok || userID == "" {
		return uuid.Nil, errors.New("unauthorized")
	}

	id, err := uuid.Parse(userID)

	if err != nil {
		return uuid.Nil, errors.New("invalid user id")
	}

	return id, nil
}

// POST /api/enrollments/courses/:courseId
func (c *EnrollmentController) Enroll(
	ctx *fiber.Ctx,
) error {

	userID, err := getUserID(ctx)

	if err != nil {
		return ctx.Status(
			fiber.StatusUnauthorized,
		).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	courseID, err := uuid.Parse(
		ctx.Params("courseId"),
	)

	if err != nil {
		return ctx.Status(
			fiber.StatusBadRequest,
		).JSON(fiber.Map{
			"success": false,
			"message": "invalid course id",
		})
	}

	result, err := c.service.Enroll(
		userID,
		courseID,
	)

	if err != nil {
		return ctx.Status(
			fiber.StatusBadRequest,
		).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(
		fiber.StatusCreated,
	).JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// GET /api/enrollments/my
func (c *EnrollmentController) GetMyEnrollments(
	ctx *fiber.Ctx,
) error {

	userID, err := getUserID(ctx)

	if err != nil {
		return ctx.Status(
			fiber.StatusUnauthorized,
		).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	result, err := c.service.GetMyEnrollments(userID)

	if err != nil {
		return ctx.Status(
			fiber.StatusInternalServerError,
		).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// GET /api/enrollments/:id
func (c *EnrollmentController) GetEnrollmentByID(
	ctx *fiber.Ctx,
) error {

	userID, err := getUserID(ctx)

	if err != nil {
		return ctx.Status(
			fiber.StatusUnauthorized,
		).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	id, err := uuid.Parse(ctx.Params("id"))

	if err != nil {
		return ctx.Status(
			fiber.StatusBadRequest,
		).JSON(fiber.Map{
			"success": false,
			"message": "invalid enrollment id",
		})
	}

	result, err := c.service.GetEnrollmentByID(
		userID,
		id,
	)

	if err != nil {
		status := fiber.StatusBadRequest

		if err.Error() == "enrollment not found" {
			status = fiber.StatusNotFound
		}

		if err.Error() == "permission denied" {
			status = fiber.StatusForbidden
		}

		return ctx.Status(status).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// DELETE /api/enrollments/:id
func (c *EnrollmentController) Unenroll(
	ctx *fiber.Ctx,
) error {

	userID, err := getUserID(ctx)

	if err != nil {
		return ctx.Status(
			fiber.StatusUnauthorized,
		).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	id, err := uuid.Parse(ctx.Params("id"))

	if err != nil {
		return ctx.Status(
			fiber.StatusBadRequest,
		).JSON(fiber.Map{
			"success": false,
			"message": "invalid enrollment id",
		})
	}

	if err := c.service.Unenroll(
		userID,
		id,
	); err != nil {

		status := fiber.StatusBadRequest

		if err.Error() == "enrollment not found" {
			status = fiber.StatusNotFound
		}

		if err.Error() == "permission denied" {
			status = fiber.StatusForbidden
		}

		return ctx.Status(status).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "unenrolled successfully",
	})
}
