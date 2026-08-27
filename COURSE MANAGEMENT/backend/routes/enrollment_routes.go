package routes

import (
	"course-management/controllers"
	"course-management/middleware"
	"course-management/repositories"
	"course-management/services"

	"github.com/gofiber/fiber/v2"
)

func SetupEnrollmentRoutes(app *fiber.App) {

	enrollmentRepository :=
		repositories.NewEnrollmentRepository()

	courseRepository :=
		repositories.NewCourseRepository()

	enrollmentService :=
		services.NewEnrollmentService(
			enrollmentRepository,
			courseRepository,
		)

	controller :=
		controllers.NewEnrollmentController(
			enrollmentService,
		)

	enrollments := app.Group(
		"/api/enrollments",
		middleware.AuthMiddleware,
	)

	enrollments.Post(
		"/courses/:courseId",
		controller.Enroll,
	)

	enrollments.Get(
		"/my",
		controller.GetMyEnrollments,
	)

	enrollments.Get(
		"/:id",
		controller.GetEnrollmentByID,
	)

	enrollments.Delete(
		"/:id",
		controller.Unenroll,
	)
}