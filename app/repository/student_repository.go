package repository

import (
	"database/sql"
	"time"
)

type StudentRepository struct {
	DB *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{DB: db}
}


// GET /students
func (r *StudentRepository) GetAll() ([]map[string]interface{}, error) {
	rows, err := r.DB.Query(`SELECT id, user_id, advisor_id FROM students`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id, userID string
		var advisorID *string

		rows.Scan(&id, &userID, &advisorID)

		result = append(result, map[string]interface{}{
			"id":         id,
			"user_id":    userID,
			"advisor_id": advisorID,
		})
	}

	return result, nil
}

// GET /students/:id
func (r *StudentRepository) GetByID(id string) (map[string]interface{}, error) {
	row := r.DB.QueryRow(`
		SELECT id, user_id, advisor_id FROM students WHERE id=$1
	`, id)

	var sid, uid string
	var aid *string
	err := row.Scan(&sid, &uid, &aid)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":         sid,
		"user_id":    uid,
		"advisor_id": aid,
	}, nil
}

func (r *StudentRepository) GetAchievementsStudents(studentID string) ([]map[string]interface{}, error) {
	rows, err := r.DB.Query(`
		SELECT id, status, created_at
		FROM achievement_references
		WHERE student_id=$1
	`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id, status string
		var created time.Time

		rows.Scan(&id, &status, &created)

		result = append(result, map[string]interface{}{
			"id":         id,
			"status":     status,
			"created_at": created,
		})
	}

	return result, nil
}

func (r *StudentRepository) UpdateAdvisor(studentID, lecturerID string) error {
	_, err := r.DB.Exec(`
		UPDATE students SET advisor_id=$1 WHERE id=$2
	`, lecturerID, studentID)

	return err
}
