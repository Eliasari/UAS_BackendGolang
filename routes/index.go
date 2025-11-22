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
) {

	// auth
	AuthRoutes(app, authService)

	// user management
	UserRoutes(app, userService, permService)

	// nanti achievement
	// AchievementRoutes(app, achievementService, permService)

	// nanti students
	// StudentRoutes(app, studentService, permService)

	// nanti lecturers
	// LecturerRoutes(app, lecturerService, permService)
}
