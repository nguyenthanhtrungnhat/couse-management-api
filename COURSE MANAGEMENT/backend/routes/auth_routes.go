package routes

import (
	"course-management/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app fiber.Router) {

	authController := controllers.NewAuthController()

	auth := app.Group("/auth")

	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
}
