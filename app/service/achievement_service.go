package service

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"

	"uas-prestasi/app/model"
	"uas-prestasi/app/repository"

	"github.com/gofiber/fiber/v2"
)

type AchievementService struct {
	MongoRepo         *repository.AchievementMongoRepository
	RefRepo           *repository.AchievementReferenceRepository
	StudentDB         *sql.DB
	PermissionService *PermissionService
	PermissionRepo    *repository.PermissionRepository
}

func NewAchievementService(mongoRepo *repository.AchievementMongoRepository,
	refRepo *repository.AchievementReferenceRepository,
	studentDB *sql.DB,
	permissionService *PermissionService,
	permissionRepo *repository.PermissionRepository) *AchievementService {

	return &AchievementService{
		MongoRepo:         mongoRepo,
		RefRepo:           refRepo,
		StudentDB:         studentDB,
		PermissionService: permissionService,
		PermissionRepo:    permissionRepo,
	}
}

// CreateDraft godoc
// @Summary Create achievement draft
// @Description Create draft prestasi.
// Details bersifat dinamis tergantung achievement_type.
// Contoh competition: { competitionName, competitionLevel, rank, medalType }
// @Tags Achievement
// @Accept json
// @Produce json
// @Param request body model.CreateAchievementRequest true "Create draft payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/achievements [post]
func (s *AchievementService) CreateDraft(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req model.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	// 1️⃣ Ambil student_id dari users
	var studentID string
	err := s.StudentDB.QueryRow(
		"SELECT id FROM students WHERE user_id = $1",
		userID,
	).Scan(&studentID)

	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "user is not a student"})
	}

	// 2️⃣ Insert ke Mongo
	achievement := &model.Achievement{
		StudentID:       studentID,
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Details:         req.Details,
		Tags:            req.Tags,
		Points:          req.Points,
	}

	mongoID, err := s.MongoRepo.InsertDraft(achievement)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to insert mongo"})
	}

	// 3️⃣ Insert ke PostgreSQL reference
	ref := &model.AchievementReference{
		StudentID:          studentID,
		MongoAchievementID: mongoID,
		Status:             "draft",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := s.RefRepo.InsertDraft(ref); err != nil {
		_ = s.MongoRepo.DeleteByID(mongoID)
		return c.Status(500).JSON(fiber.Map{"error": "failed to insert reference"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"mongo_id": mongoID,
			"status":   "draft",
		},
	})

}

// Submit godoc
// @Summary Submit achievement
// @Description Submit draft prestasi menjadi submitted
// @Tags Achievement
// @Produce json
// @Param id path string true "Achievement ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/achievements/{id}/submit [post]
func (s *AchievementService) Submit(c *fiber.Ctx) error {
	achievementID := c.Params("id")
	userID := c.Locals("user_id").(string)

	var studentID string
	err := s.StudentDB.QueryRow(
		"SELECT id FROM students WHERE user_id = $1",
		userID,
	).Scan(&studentID)

	if err != nil {
		return c.Status(403).JSON(fiber.Map{
			"error": "user is not a student",
		})
	}

	err = s.RefRepo.SubmitDraft(achievementID, studentID)
	if err == sql.ErrNoRows {
		return c.Status(400).JSON(fiber.Map{
			"error": "achievement not found or not in draft status",
		})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to submit achievement",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "achievement submitted successfully",
		"data": fiber.Map{
			"id":     achievementID,
			"status": "submitted",
		},
	})
}

// Verify godoc
// @Summary Verify achievement
// @Description Dosen wali memverifikasi prestasi mahasiswa
// @Tags Achievement
// @Produce json
// @Param id path string true "Achievement ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/achievements/{id}/verify [post]
func (s *AchievementService) Verify(c *fiber.Ctx) error {
	achievementID := c.Params("id")
	lecturerID := c.Locals("user_id").(string)

	_, err := s.RefRepo.GetForVerification(achievementID, lecturerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Prestasi tidak ditemukan atau bukan bimbingan Anda",
			})
		}

		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal memvalidasi prestasi",
		})
	}

	err = s.RefRepo.Verify(achievementID, lecturerID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal memverifikasi prestasi",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Prestasi berhasil diverifikasi",
		"data": fiber.Map{
			"id":     achievementID,
			"status": "verified",
		},
	})
}

// Reject godoc
// @Summary Reject achievement
// @Description Dosen wali menolak prestasi dengan catatan
// @Tags Achievement
// @Accept json
// @Produce json
// @Param id path string true "Achievement ID"
// @Param request body model.RejectionNote true "Rejection note"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/achievements/{id}/reject [post]
func (s *AchievementService) Reject(c *fiber.Ctx) error {
	achievementID := c.Params("id")
	lecturerID := c.Locals("user_id").(string)

	type RejectRequest struct {
		Note string `json:"note"`
	}

	var req RejectRequest
	if err := c.BodyParser(&req); err != nil || req.Note == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Catatan penolakan wajib diisi",
		})
	}

	_, err := s.RefRepo.GetForVerification(achievementID, lecturerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Prestasi tidak ditemukan atau bukan bimbingan Anda",
			})
		}

		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal memvalidasi prestasi",
		})
	}

	err = s.RefRepo.Reject(achievementID, lecturerID, req.Note)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menolak prestasi",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Prestasi berhasil ditolak",
		"data": fiber.Map{
			"id":     achievementID,
			"status": "rejected",
		},
	})
}

// List godoc
// @Summary List achievements
// @Description Menampilkan daftar prestasi sesuai role dan permission (Admin, Dosen, Mahasiswa)
// @Tags Achievement
// @Produce json
//
// @Param page query int false "Nomor halaman" default(1)
// @Param limit query int false "Jumlah data per halaman" default(10)
// @Param sort query string false "Field sorting (created_at, status)" default(created_at)
// @Param order query string false "Urutan sorting (asc | desc)" default(desc)
// @Param status query string false "Filter status prestasi (pending, approved, rejected)"
//
// @Success 200 {object} map[string]interface{} "List achievements dengan pagination"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal server error"
//
// @Security BearerAuth
// @Router /api/v1/achievements [get]
func (s *AchievementService) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	roleID := c.Locals("role_id").(string)

	// ===== Pagination =====
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// ===== Filter & Sorting =====
	sortBy := c.Query("sort", "created_at")
	order := c.Query("order", "desc")
	status := c.Query("status", "")

	var (
		refs  []model.AchievementReference
		total int
		err   error
	)

	// ===== Role-based access =====
	if ok, _ := s.PermissionService.HasPermission(roleID, "achievement:list:all"); ok {

		refs, total, err = s.RefRepo.ListAll(
			limit, offset, sortBy, order, status,
		)

	} else if ok, _ := s.PermissionService.HasPermission(roleID, "achievement:list:advisor"); ok {

		refs, total, err = s.RefRepo.ListByLecturer(
			userID, limit, offset, sortBy, order, status,
		)

	} else if ok, _ := s.PermissionService.HasPermission(roleID, "achievement:list:self"); ok {

		studentID, err2 := s.RefRepo.GetStudentIDByUser(userID)
		if err2 != nil {
			return c.Status(403).JSON(fiber.Map{
				"message": "Student data not found",
			})
		}

		refs, total, err = s.RefRepo.ListByStudent(
			studentID, limit, offset, sortBy, order, status,
		)

	} else {
		return c.Status(403).JSON(fiber.Map{
			"message": "You do not have permission to access this resource",
		})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// ===== Build response data =====
	var results []fiber.Map
	for _, ref := range refs {
		mongoData, _ := s.MongoRepo.FindByID(ref.MongoAchievementID)

		results = append(results, fiber.Map{
			"id":        ref.ID,
			"status":    ref.Status,
			"studentId": ref.StudentID,
			"mongo":     mongoData,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"meta": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": int(math.Ceil(float64(total) / float64(limit))),
		},
		"data": results,
	})
}

// Detail godoc
// @Summary Achievement detail
// @Description Detail prestasi (owner, advisor, atau admin)
// @Tags Achievement
// @Produce json
// @Param id path string true "Achievement ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/achievements/{id} [get]
func (s *AchievementService) Detail(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	// 1️⃣ Ambil reference
	ref, err := s.RefRepo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Achievement tidak ditemukan",
		})
	}

	// 2️⃣ Coba sebagai mahasiswa (owner)
	if studentID, err := s.RefRepo.GetStudentIDByUser(userID); err == nil {
		if studentID != ref.StudentID {
			// bukan owner → cek advisor
			isAdvisor, _ := s.RefRepo.IsAdvisorOf(userID, ref.StudentID)
			if !isAdvisor {
				return c.Status(403).JSON(fiber.Map{
					"error": "forbidden",
				})
			}
		}
		// owner atau advisor → lanjut
	}

	// kalau bukan mahasiswa → admin/dosen non wali → lolos karena RBAC

	// 3️⃣ Ambil detail mongo
	mongoData, err := s.MongoRepo.FindByID(ref.MongoAchievementID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "achievement detail not found",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id":                 ref.ID,
			"status":             ref.Status,
			"created_at":         ref.CreatedAt,
			"achievement_detail": mongoData,
		},
	})
}

// Update godoc
// @Summary Update achievement draft
// @Description Update draft prestasi milik mahasiswa
// @Tags Achievement
// @Accept json
// @Produce json
// @Param id path string true "Achievement ID"
// @Param request body model.UpdateAchievementRequest true "Update payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/achievements/{id} [put]
func (s *AchievementService) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	studentID, err := s.RefRepo.GetStudentIDByUser(userID)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{
			"message": "Data mahasiswa tidak ditemukan",
		})
	}

	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Payload invalid",
		})
	}

	ref, err := s.RefRepo.GetDraftByOwner(id, studentID)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{
			"message": "Draft tidak ditemukan atau bukan milik Anda",
		})
	}

	err = s.MongoRepo.UpdateByID(ref.MongoAchievementID, payload)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal update data",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Draft berhasil diupdate",
		"data": fiber.Map{
			"id":       id,
			"mongo_id": ref.MongoAchievementID,
		},
	})
}

// Delete godoc
// @Summary Delete achievement draft
// @Description Menghapus draft prestasi
// @Tags Achievement
// @Produce json
// @Param id path string true "Achievement ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/achievements/{id} [delete]
func (s *AchievementService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	studentID, err := s.RefRepo.GetStudentIDByUser(userID)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{
			"message": "Mahasiswa tidak ditemukan",
		})
	}

	ref, err := s.RefRepo.GetOwnedAchievement(id, studentID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Draft tidak ditemukan atau bukan milik Anda",
		})
	}

	if ref.Status != "draft" {
		return c.Status(409).JSON(fiber.Map{
			"message": "Hanya draft yang dapat dihapus",
		})
	}

	if err := s.RefRepo.DeleteDraft(id, studentID); err != nil {
		return c.Status(409).JSON(fiber.Map{
			"message": "Draft tidak dapat dihapus (mungkin sudah diproses)",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Draft berhasil dihapus",
		"data": fiber.Map{
			"id": id,
		},
	})
}

// UploadAttachment godoc
// @Summary Upload achievement attachment
// @Description Upload file pendukung prestasi
// @Tags Achievement
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Achievement ID"
// @Param file formData file true "Attachment file"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/achievements/{id}/attachments [post]
func (s *AchievementService) UploadAttachment(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	// ✅ Ambil student_id dari user_id
	studentID, err := s.RefRepo.GetStudentIDByUser(userID)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{
			"message": "Data mahasiswa tidak ditemukan",
		})
	}

	// ✅ Baru validasi kepemilikan achievement
	ref, err := s.RefRepo.GetOwnedAchievement(id, studentID)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{
			"message": "Achievement tidak ditemukan atau bukan milik Anda",
		})
	}

	// 🔒 Kunci kalau sudah diverifikasi
	if ref.Status == "verified" {
		return c.Status(403).JSON(fiber.Map{
			"message": "Prestasi yang sudah diverifikasi tidak bisa diubah",
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "File wajib diupload",
		})
	}

	// ✅ Path lebih aman pakai studentID asli
	uploadDir := fmt.Sprintf("uploads/%s", studentID)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal membuat folder upload",
		})
	}

	fileName := fmt.Sprintf("%s_%s", uuid.NewString(), file.Filename)
	uploadPath := fmt.Sprintf("%s/%s", uploadDir, fileName)

	if err := c.SaveFile(file, uploadPath); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal menyimpan file",
		})
	}

	attachment := model.Attachment{
		FileName:   file.Filename,
		FileURL:    uploadPath,
		FileType:   file.Header.Get("Content-Type"),
		UploadedAt: time.Now(),
	}

	err = s.MongoRepo.AddAttachment(ref.MongoAchievementID, attachment)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal menyimpan attachment ke MongoDB",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Attachment berhasil diupload",
		"file":    attachment,
	})
}

// History godoc
// @Summary Achievement history
// @Description Menampilkan riwayat status prestasi
// @Tags Achievement
// @Produce json
// @Param id path string true "Achievement ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/achievements/{id}/history [get]
func (s *AchievementService) History(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	ref, err := s.RefRepo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Achievement tidak ditemukan"})
	}

	studentID, err := s.RefRepo.GetStudentIDByUser(userID)
	if err == nil {
		isOwner := studentID == ref.StudentID
		isAdvisor, _ := s.RefRepo.IsAdvisorOf(userID, ref.StudentID)

		if !isOwner && !isAdvisor {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	history, err := s.RefRepo.GetHistory(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil history"})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"history": history,
	})
}
