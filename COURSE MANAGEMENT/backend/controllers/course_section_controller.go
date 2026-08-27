package controllers

import (
	"course-management/dto/section"
	"course-management/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CourseSectionController struct {
	service services.CourseSectionService
}

func NewCourseSectionController(
	service services.CourseSectionService,
) *CourseSectionController {

	return &CourseSectionController{
		service: service,
	}
}

func (ctrl *CourseSectionController) CreateSection(
	c *fiber.Ctx,
) error {

	var req section.CreateSectionRequest

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

	result, err := ctrl.service.CreateSection(
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

func (ctrl *CourseSectionController) GetSection(
	c *fiber.Ctx,
) error {

	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "invalid section id",
			},
		)
	}

	result, err := ctrl.service.GetSection(id)

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

func (ctrl *CourseSectionController) GetSectionsByCourse(
	c *fiber.Ctx,
) error {

	courseID, err := uuid.Parse(
		c.Params("courseId"),
	)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "invalid course id",
			},
		)
	}

	result, err := ctrl.service.GetSectionsByCourse(
		courseID,
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

func (ctrl *CourseSectionController) UpdateSection(
	c *fiber.Ctx,
) error {

	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "invalid section id",
			},
		)
	}

	var req section.UpdateSectionRequest

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

	result, err := ctrl.service.UpdateSection(
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

func (ctrl *CourseSectionController) DeleteSection(
	c *fiber.Ctx,
) error {

	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "invalid section id",
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

	if err := ctrl.service.DeleteSection(
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
			"message": "section deleted successfully",
		},
	)
}