package service

import (
	"uas-prestasi/app/model"
	"uas-prestasi/app/repository"

	"github.com/gofiber/fiber/v2"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

// GetAll godoc
// @Summary Ambil semua user
// @Description Mengambil seluruh data user
// @Tags Users
// @Produce json
// @Success 200 {object} model.GetUserResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/users [get]
func (s *UserService) GetAll(c *fiber.Ctx) error {
	users, err := s.Repo.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	if len(users) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "resource not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   users,
	})
}

// GetByID godoc
// @Summary Ambil user berdasarkan ID
// @Description Mengambil detail user berdasarkan ID
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} model.GetUserResponse
// @Failure 404 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/users/{id} [get]
func (s *UserService) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := s.Repo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}
	return c.JSON(fiber.Map{"data": user})
}

// Create godoc
// @Summary Tambah user baru
// @Description Membuat user baru
// @Tags Users
// @Accept json
// @Produce json
// @Param request body model.CreateUserRequest true "Data user"
// @Success 201 {object} model.CreateUserResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/users [post]
func (s *UserService) Create(c *fiber.Ctx) error {
	var req model.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(model.ErrorResponse{
			Error: "invalid request",
		})
	}

	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		RoleID:   req.RoleID,
	}

	if err := s.Repo.Create(&user); err != nil {
		return c.Status(500).JSON(model.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(201).JSON(model.CreateUserResponse{
		Data: user,
	})
}

// Update godoc
// @Summary Update user
// @Description Update data user berdasarkan ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body model.UpdateUserRequest true "Data user"
// @Success 200 {object} model.UpdateUserResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/users/{id} [put]
func (s *UserService) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req model.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(model.ErrorResponse{
			Error: "invalid request",
		})
	}

	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		RoleID:   req.RoleID,
		IsActive: req.IsActive,
	}

	if err := s.Repo.Update(id, &user); err != nil {
		return c.Status(500).JSON(model.ErrorResponse{
			Error: "failed to update user",
		})
	}

	updatedUser, err := s.Repo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(model.ErrorResponse{
			Error: "user not found",
		})
	}

	return c.Status(200).JSON(model.UpdateUserResponse{
		Status: "success",
		Data:   *updatedUser,
	})
}

// Delete godoc
// @Summary Hapus user
// @Description Menghapus user berdasarkan ID
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} model.MessageResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/users/{id} [delete]
func (s *UserService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := s.Repo.Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "deleted"})
}

// UpdateRole godoc
// @Summary Update role user
// @Description Mengubah role user berdasarkan ID User
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body model.UpdateRoleRequest true "Role payload"
// @Success 200 {object} model.UpdateRoleResult
// @Failure 400 {object} model.ErrorResponse
// @Failure 422 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/users/{id}/role [put]
func (s *UserService) UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")

	var req struct {
		RoleID string `json:"role_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	if req.RoleID == "" {
		return c.Status(422).JSON(fiber.Map{
			"error": "role_id is required",
		})
	}

	if err := s.Repo.UpdateRole(id, req.RoleID); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to update role",
		})
	}

	user, err := s.Repo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id":        user.ID,
			"username":  user.Username,
			"role":      user.RoleName,
			"role_id":   user.RoleID,
			"is_active": user.IsActive,
		},
	})
}
