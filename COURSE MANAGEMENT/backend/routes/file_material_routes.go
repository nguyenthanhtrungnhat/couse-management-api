package routes

import (
	"course-management/controllers"
	"course-management/middleware"

	"github.com/gofiber/fiber/v2"
)

func FileMaterialRoutes(app fiber.Router) {

	controller := controllers.NewFileMaterialController()

	materials := app.Group("/materials")

	materials.Use(middleware.AuthMiddleware)

	materials.Get(
		"/lesson/:lessonId",
		controller.GetMaterialsByLesson,
	)

	materials.Get(
		"/:id",
		controller.GetMaterialByID,
	)

	materials.Post(
		"/lesson/:lessonId",
		controller.CreateMaterial,
	)

	materials.Delete(
		"/:id",
		middleware.RequireRole("instructor"),
		controller.DeleteMaterial,
	)
}