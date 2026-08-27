package routes

import (
	"course-management/controllers"
	"course-management/middleware"

	"github.com/gofiber/fiber/v2"
)

func FileMaterialRoutes(app fiber.Router) {

	controller := controllers.NewFileMaterialController()

	materials := app.Group("/materials")

	// ========================================================
	// PUBLIC / AUTHENTICATED
	// ========================================================

	materials.Use(middleware.JWTProtected())

	// Get all materials of a lesson
	materials.Get(
		"/lesson/:lessonId",
		controller.GetMaterialsByLesson,
	)

	// Get material detail
	materials.Get(
		"/:id",
		controller.GetMaterialByID,
	)

	// ========================================================
	// INSTRUCTOR ONLY
	// ========================================================

	instructor := materials.Group(
		"",
		middleware.RequireRole("instructor"),
	)

	// Upload material
	instructor.Post(
		"/lesson/:lessonId",
		controller.CreateMaterial,
	)

	// Delete material
	instructor.Delete(
		"/:id",
		controller.DeleteMaterial,
	)
}