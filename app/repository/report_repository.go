package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReportRepository struct {
	DB    *sql.DB
	Mongo *mongo.Database
}

func NewReportRepository(db *sql.DB, mongo *mongo.Database) *ReportRepository {
	return &ReportRepository{DB: db, Mongo: mongo}
}

func (r *ReportRepository) GetStudentNamesFromStudents(ids []string) (map[string]string, error) {
	if len(ids) == 0 {
		return map[string]string{}, nil
	}

	query := `
		SELECT s.id, u.full_name
		FROM students s
		JOIN users u ON s.user_id = u.id
		WHERE s.id = ANY($1)
	`

	rows, err := r.DB.Query(query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var studentID, fullName string
		if err := rows.Scan(&studentID, &fullName); err != nil {
			return nil, err
		}
		result[studentID] = fullName
	}

	return result, nil
}

// =============================
// Ambil achievement references berdasarkan role
// =============================
func (r *ReportRepository) GetAchievementReferencesByUser(userID, roleID string) ([]string, error) {

	var roleName string
	err := r.DB.QueryRow(`SELECT name FROM roles WHERE id = $1`, roleID).Scan(&roleName)
	if err != nil {
		return nil, err
	}

	roleName = strings.ToLower(roleName)

	baseQuery := `
		SELECT ar.mongo_achievement_id
		FROM achievement_references ar
		JOIN students s ON ar.student_id = s.id
		LEFT JOIN lecturers l ON s.advisor_id = l.id
		WHERE ar.status = 'verified'
	`

	var rows *sql.Rows

	switch roleName {
	case "admin":
		rows, err = r.DB.Query(baseQuery)

	case "mahasiswa":
		rows, err = r.DB.Query(baseQuery+" AND s.user_id = $1", userID)

	case "dosen wali":
		rows, err = r.DB.Query(baseQuery+" AND l.user_id = $1", userID)

	default:
		return nil, errors.New("invalid role")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

// =============================
// Generic aggregate untuk Mongo
// =============================
func (r *ReportRepository) Aggregate(ctx context.Context, pipeline interface{}) ([]bson.M, error) {
	col := r.Mongo.Collection("achievements")

	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	// pastikan selalu return slice, walau kosong
	if result == nil {
		result = []bson.M{}
	}

	return result, nil
}

// ========================================
// ✅ Ambil Nama Mahasiswa dari SQL (Post)
// ========================================
func (r *ReportRepository) GetStudentNamesByIDs(ids []string) (map[string]string, error) {
	if len(ids) == 0 {
		return map[string]string{}, nil
	}

	query := `
		SELECT u.id, u.full_name
		FROM users u
		WHERE u.id = ANY($1)
	`

	rows, err := r.DB.Query(query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		result[id] = name
	}

	return result, nil
}
