package repository

import (
	"database/sql"
	"errors"
	"uas-prestasi/app/model"
	"time"
)

type AchievementReferenceRepository struct {
	DB *sql.DB
}

func NewAchievementReferenceRepository(db *sql.DB) *AchievementReferenceRepository {
	return &AchievementReferenceRepository{DB: db}
}

func (r *AchievementReferenceRepository) InsertDraft(ref *model.AchievementReference) error {
	query := `
		INSERT INTO achievement_references
		(student_id, mongo_achievement_id, status)
		VALUES ($1, $2, $3)
	`

	_, err := r.DB.Exec(
		query,
		ref.StudentID,
		ref.MongoAchievementID,
		ref.Status,
	)

	return err
}

func (r *AchievementReferenceRepository) SubmitDraft(id string, studentID string) error {
	query := `
		UPDATE achievement_references
		SET 
			status = 'submitted',
			submitted_at = NOW(),
			updated_at = NOW()
		WHERE id = $1 
		  AND student_id = $2
		  AND status = 'draft'
	`

	res, err := r.DB.Exec(query, id, studentID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *AchievementReferenceRepository) GetForVerification(refID string, lecturerUserID string) (*model.AchievementReference, error) {
	query := `
		SELECT ar.id, ar.student_id, ar.mongo_achievement_id, ar.status
		FROM achievement_references ar
		JOIN students s ON s.id = ar.student_id
		JOIN lecturers l ON s.advisor_id = l.id
		WHERE ar.id = $1
		  AND ar.status = 'submitted'
		  AND l.user_id = $2
	`

	var ref model.AchievementReference

	err := r.DB.QueryRow(query, refID, lecturerUserID).Scan(
		&ref.ID,
		&ref.StudentID,
		&ref.MongoAchievementID,
		&ref.Status,
	)

	if err != nil {
		return nil, err
	}

	return &ref, nil
}


func (r *AchievementReferenceRepository) Verify(id, lecturerID string) error {
	query := `
		UPDATE achievement_references
		SET 
			status = 'verified',
			verified_at = NOW(),
			verified_by = $1,
			updated_at = NOW()
		WHERE id = $2 
		  AND status = 'submitted'
	`

	res, err := r.DB.Exec(query, lecturerID, id)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("data tidak ditemukan atau belum diajukan")
	}

	return nil
}

func (r *AchievementReferenceRepository) Reject(id, lecturerID, note string) error {
	query := `
		UPDATE achievement_references
		SET 
			status = 'rejected',
			rejection_note = $1,
			verified_at = NOW(),
			verified_by = $2,
			updated_at = NOW()
		WHERE id = $3 
		  AND status = 'submitted'
	`

	res, err := r.DB.Exec(query, note, lecturerID, id)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("data tidak ditemukan atau belum diajukan")
	}

	return nil
}

func (r *AchievementReferenceRepository) ListByStudent(studentID string) ([]model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at
		FROM achievement_references
		WHERE student_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.AchievementReference

	for rows.Next() {
		var ref model.AchievementReference
		err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.SubmittedAt,
			&ref.VerifiedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, ref)
	}

	return result, nil
}

func (r *AchievementReferenceRepository) ListByLecturer(lecturerID string) ([]model.AchievementReference, error) {
	query := `
		SELECT ar.id, ar.student_id, ar.mongo_achievement_id, ar.status, ar.submitted_at, ar.verified_at
		FROM achievement_references ar
		JOIN students s ON s.id = ar.student_id
		WHERE s.advisor_id = $1
		ORDER BY ar.created_at DESC
	`

	rows, err := r.DB.Query(query, lecturerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.AchievementReference

	for rows.Next() {
		var ref model.AchievementReference
		err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.SubmittedAt,
			&ref.VerifiedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, ref)
	}

	return result, nil
}

func (r *AchievementReferenceRepository) ListAll() ([]model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at
		FROM achievement_references
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.AchievementReference

	for rows.Next() {
		var ref model.AchievementReference
		err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.SubmittedAt,
			&ref.VerifiedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, ref)
	}

	return result, nil
}

// detail achievement
func (r *AchievementReferenceRepository) GetByID(id string) (*model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, created_at
		FROM achievement_references
		WHERE id = $1
	`

	var a model.AchievementReference
	err := r.DB.QueryRow(query, id).Scan(
		&a.ID,
		&a.StudentID,
		&a.MongoAchievementID,
		&a.Status,
		&a.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (r *AchievementReferenceRepository) GetDraftByOwner(id, studentID string) (*model.AchievementReference, error) {
	query := `
		SELECT id, mongo_achievement_id, status
		FROM achievement_references
		WHERE id = $1 
		  AND student_id = $2
		  AND status = 'draft'
	`

	var ref model.AchievementReference

	err := r.DB.QueryRow(query, id, studentID).Scan(
		&ref.ID,
		&ref.MongoAchievementID,
		&ref.Status,
	)

	if err != nil {
		return nil, err
	}

	return &ref, nil
}

func (r *AchievementReferenceRepository) DeleteDraft(id, studentID string) error {
	res, err := r.DB.Exec(`
		UPDATE achievement_references
		SET status = 'deleted', updated_at = NOW()
		WHERE id = $1
		  AND student_id = $2
		  AND status = 'draft'
	`, id, studentID)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("draft not deletable")
	}

	return nil
}



func (r *AchievementReferenceRepository) GetStudentIDByUser(userID string) (string, error) {
	var studentID string

	err := r.DB.QueryRow(`
		SELECT id FROM students WHERE user_id = $1
	`, userID).Scan(&studentID)

	if err != nil {
		return "", err
	}

	return studentID, nil
}


func (r *AchievementReferenceRepository) IsAdvisorOf(userID string, studentID string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM students s
		JOIN lecturers l ON s.advisor_id = l.id
		WHERE s.id = $1
		  AND l.user_id = $2
	`

	var count int
	err := r.DB.QueryRow(query, studentID, userID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}


func (r *AchievementReferenceRepository) GetHistory(id string) ([]map[string]interface{}, error) {
	query := `
		SELECT status, created_at, submitted_at, verified_at, rejection_note, verified_by
		FROM achievement_references
		WHERE id = $1
	`

	var status string
	var createdAt, submittedAt, verifiedAt sql.NullTime
	var rejectionNote, verifiedBy sql.NullString

	err := r.DB.QueryRow(query, id).Scan(
		&status,
		&createdAt,
		&submittedAt,
		&verifiedAt,
		&rejectionNote,
		&verifiedBy,
	)
	if err != nil {
		return nil, err
	}

	history := []map[string]interface{}{}

	// Draft (selalu ada)
	if createdAt.Valid {
		history = append(history, map[string]interface{}{
			"status": "draft",
			"time":   createdAt.Time,
		})
	}

	// Submitted
	if submittedAt.Valid {
		history = append(history, map[string]interface{}{
			"status": "submitted",
			"time":   submittedAt.Time,
		})
	}

	// Verified / Rejected
	if status == "verified" && verifiedAt.Valid {
		history = append(history, map[string]interface{}{
			"status":      "verified",
			"time":        verifiedAt.Time,
			"verified_by": verifiedBy.String,
		})
	} else if status == "rejected" && verifiedAt.Valid {
		history = append(history, map[string]interface{}{
			"status": "rejected",
			"time":   verifiedAt.Time,
			"note":   rejectionNote.String,
		})
	}

	return history, nil
}

func (r *StudentRepository) GetAchievements(studentID string) ([]map[string]interface{}, error) {
	rows, err := r.DB.Query(`
		SELECT id, status, created_at 
		FROM achievement_references
		WHERE student_id=$1
	`, studentID)

	if err != nil { return nil, err }
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
	var id, status string
	var created time.Time

	if err := rows.Scan(&id, &status, &created); err != nil {
		return nil, err
	}

	result = append(result, map[string]interface{}{
		"id": id,
		"status": status,
		"created_at": created,
	})
}

	return result, nil
}

func (r *AchievementReferenceRepository) GetOwnedAchievement(id, studentID string) (*model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, created_at
		FROM achievement_references
		WHERE id = $1 AND student_id = $2
	`

	var ref model.AchievementReference
	err := r.DB.QueryRow(query, id, studentID).Scan(
		&ref.ID,
		&ref.StudentID,
		&ref.MongoAchievementID,
		&ref.Status,
		&ref.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &ref, nil
}

