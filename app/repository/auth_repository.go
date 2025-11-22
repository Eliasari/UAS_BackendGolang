package repository

import (
	"database/sql"
	"errors"
	"uas-prestasi/app/model"
)

type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (r *AuthRepository) FindByUsernameOrEmail(identifier string) (*model.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, 
		       u.full_name, u.role_id, u.is_active, r.name AS role_name
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.username = $1 OR u.email = $1
		LIMIT 1
	`

	row := r.DB.QueryRow(query, identifier)

	var user model.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
		&user.RoleName,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("not found")
	}

	return &user, err
}

func (r *AuthRepository) FindByID(id string) (*model.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.full_name, 
		       u.password_hash, u.role_id, u.is_active, 
		       r.name AS role_name
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1
		LIMIT 1
	`

	row := r.DB.QueryRow(query, id)

	var user model.User
	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.FullName,
		&user.Password, &user.RoleID, &user.IsActive,
		&user.RoleName,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
