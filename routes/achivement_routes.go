package routes

import (
	"uas-prestasi/middleware"
	"uas-prestasi/app/service"

	"github.com/gofiber/fiber/v2"
)

func AchievementRoutes(app *fiber.App, achService *service.AchievementService, permService *service.PermissionService) {
	routes := app.Group("/api/v1/achievements",
		middleware.JWTMiddleware,
	)

	routes.Post("/",
		middleware.RBAC("achievement:create", permService),
		achService.CreateDraft,
	)

	routes.Post("/:id/submit",
		middleware.RBAC("achievement:submit", permService),
		achService.Submit,
	)

	routes.Post("/:id/verify",
		middleware.RBAC("achievement:verify", permService),
		achService.Verify,
	)

	routes.Post("/:id/reject",
		middleware.RBAC("achievement:reject", permService),
		achService.Reject,
	)

	routes.Get("/:id",
	middleware.RBAC("achievement:detail", permService),
	achService.Detail,
	)

	routes.Get("/",
	middleware.RBAC("achievement:list", permService),
	achService.List,
	)

	routes.Put("/:id",
	middleware.RBAC("achievement:update", permService),
	achService.Update,
	)

	routes.Delete("/:id",
	middleware.RBAC("achievement:delete", permService),
	achService.Delete,
	)

	routes.Post("/:id/attachments",
	middleware.RBAC("achievement:upload", permService),
	achService.UploadAttachment,
	)

	routes.Get("/:id/history",
	middleware.RBAC("achievement:history", permService),
	achService.History,
	)


}
