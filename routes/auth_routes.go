package routes

import (
	"uas-prestasi/middleware"
	"uas-prestasi/app/service"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, authService *service.AuthService) {
	r := app.Group("/api/v1/auth")

	r.Post("/login", authService.Login)
	r.Post("/refresh", authService.RefreshToken)
	r.Post("/logout", middleware.JWTMiddleware, authService.Logout)
	r.Get("/profile", middleware.JWTMiddleware, authService.Profile)
}
