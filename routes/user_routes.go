package routes

import (
	"uas-prestasi/app/service"
	"uas-prestasi/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App, userService *service.UserService, permService *service.PermissionService) {
	routes := app.Group("/api/v1/users")

	routes.Get("/",
		middleware.JWTMiddleware,
		middleware.RBAC("user:manage", permService),
		userService.GetAll,
	)

	routes.Get("/:id",
		middleware.JWTMiddleware,
		middleware.RBAC("user:manage", permService),
		userService.GetByID,
	)

	routes.Post("/",
		middleware.JWTMiddleware,
		middleware.RBAC("user:manage", permService),
		userService.Create,
	)

	routes.Put("/:id",
		middleware.JWTMiddleware,
		middleware.RBAC("user:manage", permService),
		userService.Update,
	)

	routes.Delete("/:id",
		middleware.JWTMiddleware,
		middleware.RBAC("user:manage", permService),
		userService.Delete,
	)

	routes.Put("/:id/role",
		middleware.JWTMiddleware,
		middleware.RBAC("user:manage", permService),
		userService.UpdateRole,
	)
}
