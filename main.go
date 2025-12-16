package main

import (
	"os"
	"uas-prestasi/config"
	"uas-prestasi/database"
	"uas-prestasi/app/repository"
	"uas-prestasi/app/service"

	"github.com/gofiber/fiber/v2"

	_ "uas-prestasi/docs"
	"uas-prestasi/routes"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)
// @title API Prestasi Mahasiswa
// @version 1.0
// @description API untuk autentikasi dan manajemen data prestasi
// @host localhost:3000
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.LoadEnv()

	db := database.ConnectDB()
	mongoDB := database.ConnectMongo()

	app := fiber.New()

	authRepo := repository.NewAuthRepository(db)
	userRepo := repository.NewUserRepository(db)
	permRepo := repository.NewPermissionRepository(db)


	achievementRefRepo := repository.NewAchievementReferenceRepository(db)
	achievementMongoRepo := repository.NewAchievementMongoRepository(mongoDB)

	reportRepo := repository.NewReportRepository(db, mongoDB)

	studentRepo := repository.NewStudentRepository(db)
	lecturerRepo := repository.NewLecturerRepository(db)

	authService := service.NewAuthService(authRepo)
	userService := service.NewUserService(userRepo)
	permService := service.NewPermissionService(permRepo)
	reportService := service.NewReportService(reportRepo)

	achievementService := service.NewAchievementService(
		achievementMongoRepo,
		achievementRefRepo,
		db,
		permService,
		permRepo,
	)

	studentService := service.NewStudentService(studentRepo, lecturerRepo)
	lecturerService := service.NewStudentService(studentRepo, lecturerRepo)

	routes.RegisterRoutes(app, authService, userService, permService, achievementService, reportService, lecturerService, studentService)

	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	app.Listen(":" + os.Getenv("APP_PORT"))
}
