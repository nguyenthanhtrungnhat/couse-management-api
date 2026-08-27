package controllers

import (
	"course-management/dto/lesson"
	"course-management/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LessonController struct {
	service services.LessonService
}

func NewLessonController(
	service services.LessonService,
) *LessonController {
	return &LessonController{
		service: service,
	}
}

func (ctrl *LessonController) CreateLesson(
	c *fiber.Ctx,
) error {

	var req lesson.CreateLessonRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "invalid request body",
			},
		)
	}

	userID, ok := c.Locals("userID").(uuid.UUID)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"success": false,
				"message": "unauthorized",
			},
		)
	}

	result, err := ctrl.service.CreateLesson(
		userID,
		req,
	)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	return c.Status(fiber.StatusCreated).JSON(
		fiber.Map{
			"success": true,
			"data":    result,
		},
	)
}

func (ctrl *LessonController) GetLesson(
	c *fiber.Ctx,
) error {

	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "invalid lesson id",
			},
		)
	}

	result, err := ctrl.service.GetLesson(id)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	return c.JSON(
		fiber.Map{
			"success": true,
			"data":    result,
		},
	)
}

func (ctrl *LessonController) GetLessonsBySection(
	c *fiber.Ctx,
) error {

	sectionID, err := uuid.Parse(
		c.Params("sectionId"),
	)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "invalid section id",
			},
		)
	}

	result, err := ctrl.service.GetLessonsBySection(
		sectionID,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	return c.JSON(
		fiber.Map{
			"success": true,
			"data":    result,
		},
	)
}

func (ctrl *LessonController) UpdateLesson(
	c *fiber.Ctx,
) error {

	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "invalid lesson id",
			},
		)
	}

	var req lesson.UpdateLessonRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "invalid request body",
			},
		)
	}

	userID, ok := c.Locals("userID").(uuid.UUID)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"success": false,
				"message": "unauthorized",
			},
		)
	}

	result, err := ctrl.service.UpdateLesson(
		userID,
		id,
		req,
	)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	return c.JSON(
		fiber.Map{
			"success": true,
			"data":    result,
		},
	)
}

func (ctrl *LessonController) DeleteLesson(
	c *fiber.Ctx,
) error {

	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "invalid lesson id",
			},
		)
	}

	userID, ok := c.Locals("userID").(uuid.UUID)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"success": false,
				"message": "unauthorized",
			},
		)
	}

	if err := ctrl.service.DeleteLesson(
		userID,
		id,
	); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	return c.JSON(
		fiber.Map{
			"success": true,
			"message": "lesson deleted successfully",
		},
	)
}