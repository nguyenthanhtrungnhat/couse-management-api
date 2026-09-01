package main

import (
	"log"

	"course-management/config"
	"course-management/routes"
	//"course-management/seeders"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("❌ Cannot load .env")
	}

	// Connect PostgreSQL
	config.ConnectDatabase()

	// Seed development data
	//seeders.SeedAll()

	// Create Fiber application
	app := fiber.New()

	// =========================
	// Routes
	// =========================

	// Authentication
	routes.AuthRoutes(app)

	// Course
	routes.SetupCourseRoutes(app)

	// Course Section
	routes.SetupCourseSectionRoutes(app)

	// Lesson
	routes.SetupLessonRoutes(app)

	// File Material
	routes.FileMaterialRoutes(app)
	routes.SetupEnrollmentRoutes(app)

	// =========================
	// Health Check
	// =========================

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Course Management API",
		})
	})

	// =========================
	// Start Server
	// =========================

	log.Println("🚀 Server running on http://localhost:3000")

	log.Fatal(
		app.Listen(":3000"),
	)
}
