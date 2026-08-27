package routes

import (
	"course-management/controllers"
	"course-management/middleware"
	"course-management/repositories"
	"course-management/services"

	"github.com/gofiber/fiber/v2"
)

func SetupLessonRoutes(
	app *fiber.App,
) {

	lessonRepository :=
		repositories.NewLessonRepository()

	sectionRepository :=
		repositories.NewCourseSectionRepository()

	courseRepository :=
		repositories.NewCourseRepository()

	lessonService :=
		services.NewLessonService(
			lessonRepository,
			sectionRepository,
			courseRepository,
		)

	lessonController :=
		controllers.NewLessonController(
			lessonService,
		)

	lessons := app.Group("/api/lessons")

	// Public
	lessons.Get(
		"/section/:sectionId",
		lessonController.GetLessonsBySection,
	)

	lessons.Get(
		"/:id",
		lessonController.GetLesson,
	)

	// Authenticated
	authenticated := lessons.Group(
		"/",
		middleware.AuthMiddleware,
	)

	authenticated.Post(
		"/",
		lessonController.CreateLesson,
	)

	authenticated.Put(
		"/:id",
		lessonController.UpdateLesson,
	)

	authenticated.Delete(
		"/:id",
		lessonController.DeleteLesson,
	)
}