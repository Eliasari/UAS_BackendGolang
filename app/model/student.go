package model

import "time"

type Student struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	AdvisorID *string `json:"advisor_id"`
}

type AdviseeResponse struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	AdvisorID *string `json:"advisor_id"`
}

type AssignAdvisorResponse struct {
	Message   string `json:"message"`
	AdvisorID string `json:"advisor_id"`
}

type StudentAchievementResponse struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateAdvisorRequest struct {
	AdvisorID string `json:"advisor_id"`
}