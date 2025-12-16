package model

type Student struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	AdvisorID *string `json:"advisor_id"`
}

type UpdateAdvisorRequest struct {
	AdvisorID string `json:"advisor_id"`
}