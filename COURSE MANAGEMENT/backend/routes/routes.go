package routes

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api")

	AuthRoutes(api)
	UserRoutes(api)
	
	// CourseRoutes(api)

	// PaymentRoutes(api)

	// CommentRoutes(api)
}