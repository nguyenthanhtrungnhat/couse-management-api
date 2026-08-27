package controllers

import (
	"course-management/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FileMaterialController struct {
	service services.FileMaterialService
}

func NewFileMaterialController() *FileMaterialController {
	return &FileMaterialController{
		service: services.NewFileMaterialService(),
	}
}

// ============================================================
// CREATE / UPLOAD MATERIAL
// POST /materials/lesson/:lessonId
// ============================================================

func (c *FileMaterialController) CreateMaterial(ctx *fiber.Ctx) error {

	// Get instructor ID from JWT
	userID, ok := ctx.Locals("user_id").(string)

	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized",
		})
	}

	instructorID, err := uuid.Parse(userID)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "invalid user id",
		})
	}

	// Get lesson ID
	lessonID, err := uuid.Parse(ctx.Params("lessonId"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid lesson id",
		})
	}

	// Get uploaded file
	file, err := ctx.FormFile("file")

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "file is required",
		})
	}

	// Create material
	material, err := c.service.CreateMaterial(
		instructorID,
		lessonID,
		file,
	)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "material uploaded successfully",
		"data":    material,
	})
}

// ============================================================
// GET MATERIALS BY LESSON
// GET /materials/lesson/:lessonId
// ============================================================

func (c *FileMaterialController) GetMaterialsByLesson(
	ctx *fiber.Ctx,
) error {

	lessonID, err := uuid.Parse(ctx.Params("lessonId"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid lesson id",
		})
	}

	materials, err := c.service.GetMaterialsByLesson(lessonID)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    materials,
	})
}

// ============================================================
// GET MATERIAL BY ID
// GET /materials/:id
// ============================================================

func (c *FileMaterialController) GetMaterialByID(
	ctx *fiber.Ctx,
) error {

	materialID, err := uuid.Parse(ctx.Params("id"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid material id",
		})
	}

	material, err := c.service.GetMaterialByID(materialID)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    material,
	})
}

// ============================================================
// DELETE MATERIAL
// DELETE /materials/:id
// ============================================================

func (c *FileMaterialController) DeleteMaterial(
	ctx *fiber.Ctx,
) error {

	// Get instructor ID from JWT
	userID, ok := ctx.Locals("user_id").(string)

	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized",
		})
	}

	instructorID, err := uuid.Parse(userID)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "invalid user id",
		})
	}

	// Parse material ID
	materialID, err := uuid.Parse(ctx.Params("id"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid material id",
		})
	}

	// Delete material
	err = c.service.DeleteMaterial(
		instructorID,
		materialID,
	)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "material deleted successfully",
	})
}
