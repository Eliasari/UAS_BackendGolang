package service

import "github.com/gofiber/fiber/v2"

type AchievementService struct {}

func NewAchievementService() *AchievementService {
	return &AchievementService{}
}

func (s *AchievementService) CreateAchievement(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "CreateAchievement belum diimplementasi",
	})
}
