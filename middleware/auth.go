package middleware

import (
	"uas-prestasi/utils"

	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")

	if len(tokenString) < 8 || tokenString[:7] != "Bearer " {
		return c.Status(401).JSON(fiber.Map{"error": "missing or invalid token"})
	}

	tokenString = tokenString[7:]

	claims, err := utils.ParseToken(tokenString)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
	}

	c.Locals("user_id", claims["user_id"])
	c.Locals("role_id", claims["role_id"])

	return c.Next()
}
