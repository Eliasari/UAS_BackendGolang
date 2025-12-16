package model

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
