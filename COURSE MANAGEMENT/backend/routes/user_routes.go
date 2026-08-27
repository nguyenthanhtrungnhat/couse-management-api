package routes

import (
	"course-management/controllers"
	"course-management/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app fiber.Router) {

	controller := controllers.NewUserController()

	users := app.Group("/users")

	users.Use(middleware.AuthMiddleware)

	users.Get("/profile", controller.Profile)

	users.Put("/profile", controller.UpdateProfile)
}
