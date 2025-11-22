package routes

import (
	"uas-prestasi/app/service"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, authService *service.AuthService) {
	r := app.Group("/api/v1/auth")

	r.Post("/login", authService.Login)
	r.Post("/refresh", authService.RefreshToken)
	r.Post("/logout", authService.Logout)
	r.Get("/profile", authService.Profile)
}
