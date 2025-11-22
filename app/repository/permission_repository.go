package repository

import (
	"database/sql"
)

type PermissionRepository struct {
	DB *sql.DB
}

func NewPermissionRepository(db *sql.DB) *PermissionRepository {
	return &PermissionRepository{DB: db}
}

func (r *PermissionRepository) GetPermissionsByRole(roleID string) ([]string, error) {
	const q = `
		SELECT p.name
		FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = $1
	`
	rows, err := r.DB.Query(q, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		perms = append(perms, name)
	}
	return perms, rows.Err()
}

// Convenience: ambil semua permission untuk banyak role
func (r *PermissionRepository) GetPermissionsByRoles(roleIDs []string) ([]string, error) {
	// implementasi per-user multi-role
	return nil, nil
}
