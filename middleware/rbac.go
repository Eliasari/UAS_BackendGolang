package middleware

import (
	"uas-prestasi/app/service"

	"github.com/gofiber/fiber/v2"
)

func RBAC(requiredPermission string, permService *service.PermissionService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		roleIDRaw := c.Locals("role_id")

		if roleIDRaw == nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "role_id missing",
			})
		}
		var roleID string

		switch v := roleIDRaw.(type) {
		case string:
			roleID = v
		case []byte:
			roleID = string(v)
		default:
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid role_id type",
			})
		}

		has, err := permService.HasPermission(roleID, requiredPermission)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "permission check failed",
			})
		}

		if !has {
			return c.Status(403).JSON(fiber.Map{
				"error": "forbidden",
				"required_permission": "You do not have permission to perform this action: " + requiredPermission,
			})
		}

		return c.Next()
	}
}
