package model

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	User         User   `json:"user"`
}

type LoginUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Role     string `json:"role"`
	IsActive bool   `json:"isActive"`
}

type LoginDataResponse struct {
	Token        string            `json:"token"`
	RefreshToken string            `json:"refreshToken"`
	User         LoginUserResponse `json:"user"`
}

type ProfileResponseData struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Role     string `json:"role"`
	IsActive bool   `json:"isActive"`
}

type ProfileResponse struct {
	Status string              `json:"status"`
	Data   ProfileResponseData `json:"data"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenData struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenResponse struct {
	Status string            `json:"status"`
	Data   RefreshTokenData  `json:"data"`
}
