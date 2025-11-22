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

// GET /users
func (s *UserService) GetAll(c *fiber.Ctx) error {
	users, err := s.Repo.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": users})
}

// GET /users/:id
func (s *UserService) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := s.Repo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}
	return c.JSON(fiber.Map{"data": user})
}

// POST /users
func (s *UserService) Create(c *fiber.Ctx) error {
	var u model.User
	if err := c.BodyParser(&u); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	if err := s.Repo.Create(&u); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": u})
}

// PUT /users/:id
func (s *UserService) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var u model.User

	if err := c.BodyParser(&u); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	if err := s.Repo.Update(id, &u); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "updated"})
}

// DELETE /users/:id
func (s *UserService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := s.Repo.Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "deleted"})
}

// PUT /users/:id/role
func (s *UserService) UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")

	var req struct {
		RoleID string `json:"role_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	if err := s.Repo.UpdateRole(id, req.RoleID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "role updated"})
}
