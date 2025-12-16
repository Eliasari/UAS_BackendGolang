package service

import (
	"context"

	"uas-prestasi/app/model"
	"uas-prestasi/app/repository"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReportService struct {
	Repo *repository.ReportRepository
}

func NewReportService(repo *repository.ReportRepository) *ReportService {
	return &ReportService{Repo: repo}
}

// Statistics godoc
// @Summary Achievement statistics report
// @Tags Report
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.StatisticsResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/reports/statistics [get]
func (s *ReportService) Statistics(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	roleID := c.Locals("role_id").(string)

	ids, err := s.Repo.GetAchievementReferencesByUser(userID, roleID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch references"})
	}

	if len(ids) == 0 {
		return c.JSON(model.StatisticsResponse{})
	}

	objIDs := []interface{}{}
	for _, id := range ids {
		if oid, err := primitive.ObjectIDFromHex(id); err == nil {
			objIDs = append(objIDs, oid)
		} else {
			objIDs = append(objIDs, id)
		}
	}

	ctx := context.Background()

	typePipeline := bson.A{
		bson.M{"$match": bson.M{"_id": bson.M{"$in": objIDs}}},
		bson.M{"$group": bson.M{
			"_id":   bson.M{"$ifNull": []interface{}{"$achievementType", "unknown"}},
			"total": bson.M{"$sum": 1},
		}},
	}

	periodPipeline := bson.A{
		bson.M{"$match": bson.M{"_id": bson.M{"$in": objIDs}}},
		bson.M{"$group": bson.M{
			"_id":   bson.M{"$dateToString": bson.M{"format": "%Y-%m", "date": "$createdAt"}},
			"total": bson.M{"$sum": 1},
		}},
	}

	topPipeline := bson.A{
		bson.M{"$match": bson.M{"_id": bson.M{"$in": objIDs}}},
		bson.M{"$group": bson.M{
			"_id":   "$studentId",
			"total": bson.M{"$sum": 1},
		}},
		bson.M{"$sort": bson.M{"total": -1}},
		bson.M{"$limit": 5},
	}

	levelPipeline := bson.A{
		bson.M{"$match": bson.M{"_id": bson.M{"$in": objIDs}}},
		bson.M{"$group": bson.M{
			"_id":   bson.M{"$ifNull": []interface{}{"$details.competitionLevel", "unknown"}},
			"total": bson.M{"$sum": 1},
		}},
	}

	typeAgg, _ := s.Repo.Aggregate(ctx, typePipeline)
	periodAgg, _ := s.Repo.Aggregate(ctx, periodPipeline)
	topAgg, _ := s.Repo.Aggregate(ctx, topPipeline)
	levelAgg, _ := s.Repo.Aggregate(ctx, levelPipeline)

	// =====================
	// BUILD RESPONSE
	// =====================

	resp := model.StatisticsResponse{}

	for _, t := range typeAgg {
		code := t["_id"].(string)
		resp.TotalPerType = append(resp.TotalPerType, model.StatisticItem{
			Code:  code,
			Name:  model.AchievementTypeMap[code],
			Total: int(t["total"].(int32)),
		})
	}

	for _, p := range periodAgg {
		resp.TotalPerPeriod = append(resp.TotalPerPeriod, model.StatisticItem{
			Code:  p["_id"].(string),
			Name:  p["_id"].(string),
			Total: int(p["total"].(int32)),
		})
	}

	studentIDs := []string{}
	for _, t := range topAgg {
		studentIDs = append(studentIDs, t["_id"].(string))
	}

	nameMap, _ := s.Repo.GetStudentNamesFromStudents(studentIDs)
	for _, t := range topAgg {
		id := t["_id"].(string)
		resp.TopStudents = append(resp.TopStudents, model.StatisticTopStudentItem{
			StudentID:   id,
			StudentName: nameMap[id],
			Total:       int(t["total"].(int32)),
		})
	}

	for _, lv := range levelAgg {
		code := lv["_id"].(string)
		resp.CompetitionLevel = append(resp.CompetitionLevel, model.StatisticItem{
			Code:  code,
			Name:  model.CompetitionLevelMap[code],
			Total: int(lv["total"].(int32)),
		})
	}

	return c.JSON(resp)
}

// StudentReport godoc
// @Summary Student achievement report
// @Tags Report
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Success 200 {array} model.StudentAchievementReportResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/reports/student/{id} [get]
func (s *ReportService) StudentReport(c *fiber.Ctx) error {
	studentID := c.Params("id")
	ctx := context.Background()

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"studentId": studentID}}},
	}

	data, err := s.Repo.Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch report"})
	}

	return c.JSON(data)
}
