package repository

import (
	"database/sql"
	"uas-prestasi/app/model"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetAll() ([]model.User, error) {
	rows, err := r.DB.Query(`
		SELECT u.id, u.username, u.email, u.full_name, u.role_id, r.name AS role_name, 
		       u.is_active
		FROM users u
		LEFT JOIN roles r ON r.id = u.role_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(
			&u.ID, &u.Username, &u.Email, &u.FullName,
			&u.RoleID, &u.RoleName, &u.IsActive,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) GetByID(id string) (*model.User, error) {
	var u model.User

	err := r.DB.QueryRow(`
		SELECT u.id, u.username, u.email, u.full_name, 
		       u.role_id, r.name AS role_name, u.is_active
		FROM users u
		LEFT JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1
	`, id).Scan(
		&u.ID, &u.Username, &u.Email, &u.FullName,
		&u.RoleID, &u.RoleName, &u.IsActive,
	)

	if err != nil {
		return nil, err
	}
	return &u, nil
}


func (r *UserRepository) Create(u *model.User) error {
	err := r.DB.QueryRow(`
		INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active)
		VALUES (uuid_generate_v4(), $1, $2, crypt($3, gen_salt('bf')), $4, $5, true)
		RETURNING id
	`, u.Username, u.Email, u.Password, u.FullName, u.RoleID).Scan(&u.ID)

	return err
}


func (r *UserRepository) Update(id string, u *model.User) error {
	_, err := r.DB.Exec(`
		UPDATE users SET 
			username=$1, email=$2, full_name=$3, role_id=$4, is_active=$5, updated_at=NOW()
		WHERE id=$6
	`,
		u.Username, u.Email, u.FullName, u.RoleID, u.IsActive, id,
	)
	return err
}


func (r *UserRepository) Delete(id string) error {
	_, err := r.DB.Exec(`DELETE FROM users WHERE id=$1`, id)
	return err
}


func (r *UserRepository) UpdateRole(id, roleID string) error {
	_, err := r.DB.Exec(`
		UPDATE users SET role_id=$1, updated_at=NOW()
		WHERE id=$2
	`, roleID, id)
	return err
}

