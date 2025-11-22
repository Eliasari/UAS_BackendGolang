package main

import (
	"os"
	"uas-prestasi/config"
	"uas-prestasi/database"
	"uas-prestasi/app/repository"
	"uas-prestasi/app/service"
	"uas-prestasi/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()

	db := database.ConnectDB()

	app := fiber.New()

	authRepo := repository.NewAuthRepository(db)
	userRepo := repository.NewUserRepository(db)
	permRepo := repository.NewPermissionRepository(db)

	authService := service.NewAuthService(authRepo)
	userService := service.NewUserService(userRepo)
	permService := service.NewPermissionService(permRepo)

	routes.RegisterRoutes(app, authService, userService, permService)


	app.Listen(":" + os.Getenv("APP_PORT"))
}
