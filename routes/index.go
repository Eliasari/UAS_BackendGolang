package routes

import (
	"uas-prestasi/app/service"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(
	app *fiber.App,
	authService *service.AuthService,
	userService *service.UserService,
	permService *service.PermissionService,
	achievementService *service.AchievementService,
	reportService *service.ReportService,
	lecturerService *service.StudentService,
	studentService *service.StudentService,
) {

	// auth
	AuthRoutes(app, authService)

	// user management
	UserRoutes(app, userService, permService)

	// nanti achievement
	AchievementRoutes(app, achievementService, permService)

	// nanti students
	StudentLecturerRoutes(app, studentService, permService)

	// nanti lecturers
	StudentLecturerRoutes(app, lecturerService, permService)

	// reports
	ReportRoutes(app, reportService, permService)
}
