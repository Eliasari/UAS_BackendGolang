package routes

import (
	"uas-prestasi/middleware"
	"uas-prestasi/app/service"

	"github.com/gofiber/fiber/v2"
)

func StudentLecturerRoutes(
	app *fiber.App,
	studentService *service.StudentService,
	permService *service.PermissionService,
) {
	routes := app.Group("/api/v1")

	// Admin, Dosen
	routes.Get("/students",
		middleware.JWTMiddleware,
		middleware.RBAC("student:read", permService),
		studentService.GetStudents,
	)

	// Admin, Dosen
	routes.Get("/students/:id",
		middleware.JWTMiddleware,
		middleware.RBAC("student:read", permService),
		studentService.GetStudent,
	)

	// Mahasiswa (own), Dosen (advisee), Admin
	routes.Get("/students/:id/achievements",
		middleware.JWTMiddleware,
		middleware.RBAC("student:list", permService),
		studentService.GetStudentAchievements,
	)

	// Admin only → assign dosen wali
	routes.Put("/students/:id/advisor",
		middleware.JWTMiddleware,
		middleware.RBAC("student:set-advisor", permService),
		studentService.AssignAdvisor,
	)

	// Admin
	routes.Get("/lecturers",
		middleware.JWTMiddleware,
		middleware.RBAC("lecturer:list", permService),
		studentService.GetLecturers,
	)

	// Admin, Dosen
	routes.Get("/lecturers/:id/advisees",
		middleware.JWTMiddleware,
		middleware.RBAC("lecturer:list", permService),
		studentService.GetAdvisees,
	)
}
