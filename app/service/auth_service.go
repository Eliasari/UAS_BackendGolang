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

// Login godoc
// @Summary Login user
// @Description Login menggunakan username/email dan password, menghasilkan JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login payload"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 403 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/auth/login [post]
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

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate token baru dari refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.RefreshTokenRequest true "Refresh token payload"
// @Param request body model.RefreshTokenRequest true "Refresh token payload"
// @Success 200 {object} model.RefreshTokenResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /api/v1/auth/refresh [post]
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

// Logout godoc
// @Summary Logout user
// @Description Logout user (client-side token invalidation)
// @Tags Auth
// @Produce json
// @Success 200 {object} model.MessageResponse
// @Security BearerAuth
// @Router /api/v1/auth/logout [post]
func (s *AuthService) Logout(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"message": "logged out",
	})
}

// Profile godoc
// @Summary Get user profile
// @Description Ambil data user berdasarkan token JWT
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.ProfileResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /api/v1/auth/profile [get]
func (s *AuthService) Profile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

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

