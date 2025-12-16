package routes

import (
	"uas-prestasi/middleware"
	"uas-prestasi/app/service"

	"github.com/gofiber/fiber/v2"
)

func ReportRoutes(app *fiber.App, reportService *service.ReportService, permService *service.PermissionService) {

	routes := app.Group("/api/v1/reports")

	routes.Get("/statistics",
		middleware.JWTMiddleware,
		middleware.RBAC("report:view", permService),
		reportService.Statistics,
	)

	routes.Get("/student/:id",
		middleware.JWTMiddleware,
		middleware.RBAC("report:view", permService),
		reportService.StudentReport,
	)
}
