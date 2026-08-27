package routes

import (
	"course-management/controllers"
	"course-management/middleware"
	"course-management/repositories"
	"course-management/services"

	"github.com/gofiber/fiber/v2"
)

func SetupCourseSectionRoutes(
	app *fiber.App,
) {

	sectionRepository :=
		repositories.NewCourseSectionRepository()

	courseRepository :=
		repositories.NewCourseRepository()

	sectionService :=
		services.NewCourseSectionService(
			sectionRepository,
			courseRepository,
		)

	sectionController :=
		controllers.NewCourseSectionController(
			sectionService,
		)

	sections := app.Group(
		"/api/sections",
	)

	sections.Get(
		"/course/:courseId",
		sectionController.GetSectionsByCourse,
	)

	sections.Get(
		"/:id",
		sectionController.GetSection,
	)

	authenticated := sections.Group(
		"/",
		middleware.AuthMiddleware,
	)

	authenticated.Post(
		"/",
		sectionController.CreateSection,
	)

	authenticated.Put(
		"/:id",
		sectionController.UpdateSection,
	)

	authenticated.Delete(
		"/:id",
		sectionController.DeleteSection,
	)
}