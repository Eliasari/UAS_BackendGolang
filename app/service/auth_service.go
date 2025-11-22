package service

import (
	"uas-prestasi/app/model"
	"uas-prestasi/app/repository"
	"uas-prestasi/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo *repository.AuthRepository
}

func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{Repo: repo}
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	user, err := s.Repo.FindByUsernameOrEmail(req.Username)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	if !user.IsActive {
		return c.Status(403).JSON(fiber.Map{
			"error": "user inactive",
		})
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid credentials",
		})
	}

	token, err := utils.GenerateToken(user.ID, user.RoleID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to generate token",
		})
	}

	refresh, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to generate refresh token",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token":        token,
			"refreshToken": refresh,
			"user": fiber.Map{
				"id":         user.ID,
				"username":   user.Username,
				"email":      user.Email,
				"fullName":   user.FullName,
				"role":       user.RoleName,
				"isActive":   user.IsActive,
			},
		},
	})
}

func (s *AuthService) RefreshToken(c *fiber.Ctx) error {
	var body struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "bad request"})
	}

	claims, err := utils.ParseToken(body.RefreshToken)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid refresh token"})
	}

	userID := claims["user_id"].(string)

	// ambil user dari DB
	user, err := s.Repo.FindByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}

	// generate token baru
	newToken, _ := utils.GenerateToken(user.ID, user.RoleID)
	newRefresh, _ := utils.GenerateRefreshToken(user.ID)

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token":        newToken,
			"refreshToken": newRefresh,
		},
	})
}

func (s *AuthService) Logout(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"message": "logged out",
	})
}

func (s *AuthService) Profile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	user, err := s.Repo.FindByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"fullName":  user.FullName,
			"role":      user.RoleName,
			"isActive":  user.IsActive,
		},
	})
}

