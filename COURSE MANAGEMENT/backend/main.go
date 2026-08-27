package main

import (
	"log"

	"course-management/config"
	"course-management/routes"
	"course-management/seeders"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("❌ Cannot load .env")
	}

	config.ConnectDatabase()

	seeders.SeedRoles()

	app := fiber.New()

	// Register Auth route
	routes.AuthRoutes(app)

	// Existing routes
	routes.SetupCourseRoutes(app)
	routes.SetupCourseSectionRoutes(app)
	routes.SetupLessonRoutes(app)

	// Health check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Course Management API",
		})
	})

	log.Println("🚀 Server running on http://localhost:3000")

	log.Fatal(
		app.Listen(":3000"),
	)
}