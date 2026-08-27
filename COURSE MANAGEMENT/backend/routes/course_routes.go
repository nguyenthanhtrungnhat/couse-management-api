package routes

import (
	"course-management/controllers"
	"course-management/middleware"
	"course-management/repositories"
	"course-management/services"

	"github.com/gofiber/fiber/v2"
)

func SetupCourseRoutes(app *fiber.App) {

	courseRepository := repositories.NewCourseRepository()
	categoryRepository := repositories.NewCategoryRepository()

	courseService := services.NewCourseService(
		courseRepository,
		categoryRepository,
	)

	courseController := controllers.NewCourseController(
		courseService,
	)

	courses := app.Group("/api/courses")

	// =========================
	// PUBLIC ROUTES
	// =========================

	courses.Get(
		"/",
		courseController.SearchCourses,
	)

	courses.Get(
		"/published",
		courseController.GetPublishedCourses,
	)

	courses.Get(
		"/slug/:slug",
		courseController.GetCourseBySlug,
	)

	// =========================
	// AUTHENTICATED ROUTES
	// =========================

	authenticated := courses.Group(
		"/",
		middleware.AuthMiddleware,
	)

	authenticated.Get(
		"/my",
		courseController.GetMyCourses,
	)

	authenticated.Post(
		"/",
		courseController.CreateCourse,
	)

	authenticated.Put(
		"/:id",
		courseController.UpdateCourse,
	)

	authenticated.Delete(
		"/:id",
		courseController.DeleteCourse,
	)

	// =========================
	// PUBLIC COURSE BY ID
	// =========================

	courses.Get(
		"/:id",
		courseController.GetCourseByID,
	)
}
