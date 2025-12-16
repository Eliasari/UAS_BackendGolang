// app/repository/lecturer_repository.go
package repository

import "database/sql"

type LecturerRepository struct {
	DB *sql.DB
}

func NewLecturerRepository(db *sql.DB) *LecturerRepository {
	return &LecturerRepository{DB: db}
}


// GET /lecturers
func (r *LecturerRepository) GetAll() ([]map[string]interface{}, error) {
	rows, err := r.DB.Query(`
		SELECT l.id, l.user_id, u.username
		FROM lecturers l
		JOIN users u ON u.id = l.user_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id, userID, name string
		rows.Scan(&id, &userID, &name)

		result = append(result, map[string]interface{}{
			"id":      id,
			"user_id": userID,
			"name":    name,
		})
	}

	return result, nil
}

// GET /lecturers/:id/advisees
func (r *LecturerRepository) GetAdvisees(lecturerID string) ([]map[string]interface{}, error) {
	rows, err := r.DB.Query(`
		SELECT id, user_id, advisor_id
		FROM students
		WHERE advisor_id=$1
	`, lecturerID)
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
