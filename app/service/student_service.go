package service

import (
	"uas-prestasi/app/repository"

	"github.com/gofiber/fiber/v2"
)

type StudentService struct {
	StudentRepo  *repository.StudentRepository
	LecturerRepo *repository.LecturerRepository
}

func NewStudentService(sr *repository.StudentRepository, lr *repository.LecturerRepository) *StudentService {
	return &StudentService{sr, lr}
}

// GetStudents godoc
// @Summary Get all students
// @Description Menampilkan daftar seluruh mahasiswa
// @Tags Student
// @Produce json
// @Success 200 {array} model.Student
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/students [get]
func (s *StudentService) GetStudents(c *fiber.Ctx) error {
	data, err := s.StudentRepo.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}

// GetStudent godoc
// @Summary Get student detail
// @Description Menampilkan detail mahasiswa berdasarkan ID
// @Tags Student
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} model.Student
// @Failure 404 {object} model.MessageResponse
// @Security BearerAuth
// @Router /api/v1/students/{id} [get]
func (s *StudentService) GetStudent(c *fiber.Ctx) error {
	id := c.Params("id")

	data, err := s.StudentRepo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "student not found"})
	}
	return c.JSON(data)
}

// GetStudentAchievements godoc
// @Summary Get student achievements
// @Description Menampilkan daftar prestasi milik mahasiswa
// @Tags Student
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {array} model.StudentAchievementResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/students/{id}/achievements [get]
func (s *StudentService) GetStudentAchievements(c *fiber.Ctx) error {
	id := c.Params("id")

	data, err := s.StudentRepo.GetAchievements(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(data)
}

// AssignAdvisor godoc
// @Summary Assign academic advisor
// @Description Menentukan dosen wali untuk mahasiswa
// @Tags Student
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param request body model.UpdateAdvisorRequest true "Advisor payload"
// @Success 200 {object} model.AssignAdvisorResponse
// @Failure 400 {object} model.MessageResponse
// @Failure 500 {object} model.MessageResponse
// @Security BearerAuth
// @Router /api/v1/students/{id}/advisor [put]
func (s *StudentService) AssignAdvisor(c *fiber.Ctx) error {
	studentID := c.Params("id")

	var body struct {
		AdvisorID string `json:"advisor_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid payload"})
	}

	err := s.StudentRepo.UpdateAdvisor(studentID, body.AdvisorID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "gagal assign advisor"})
	}

	return c.JSON(fiber.Map{"message": "advisor assigned", "advisor_id": body.AdvisorID})
}

// GetLecturers godoc
// @Summary Get all lecturers
// @Description Menampilkan daftar dosen
// @Tags Lecturer
// @Produce json
// @Success 200 {array} model.AdviseeResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/lecturers [get]
func (s *StudentService) GetLecturers(c *fiber.Ctx) error {
	data, err := s.LecturerRepo.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}

// GetAdvisees godoc
// @Summary Get advisees
// @Description Menampilkan daftar mahasiswa bimbingan dosen
// @Tags Lecturer
// @Produce json
// @Param id path string true "Lecturer ID"
// @Success 200 {array} model.AdviseeResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/lecturers/{id}/advisees [get]
func (s *StudentService) GetAdvisees(c *fiber.Ctx) error {
	id := c.Params("id")

	data, err := s.LecturerRepo.GetAdvisees(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(data)
}
