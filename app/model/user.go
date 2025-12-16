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

type GetUserResponse struct {
	Status string `json:"status"`
	Data   []User `json:"data"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	RoleID   string `json:"role_id"`
}

type CreateUserResponse struct {
	Status string `json:"status"`
	Data   User   `json:"data"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	RoleID   string `json:"role_id"`
	IsActive bool   `json:"is_active"`
}

type UpdateUserResponse struct {
	Status string `json:"status"`
	Data   User   `json:"data"`
}

type UpdateRoleRequest struct {
	RoleID string `json:"role_id"`
}

type UpdateRoleResult struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	RoleID   string `json:"role_id"`
	IsActive bool   `json:"is_active"`
}



