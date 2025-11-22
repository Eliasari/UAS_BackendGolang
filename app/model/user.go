package model

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	FullName  string `json:"full_name"`
	RoleID    string `json:"role_id"`
	RoleName  string `json:"role_name,omitempty"`
	IsActive  bool   `json:"is_active"`
}
