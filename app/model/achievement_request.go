package model

import "time"

type CreateAchievementRequest struct {
	AchievementType string                 `json:"achievement_type"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Details         map[string]interface{} `json:"details"`
	Points          int                    `json:"points" example:"50"`
	Tags            []string               `json:"tags"`
}

type CreateDraftResponse struct {
	Status string                `json:"status"`
	Data   CreateDraftDataResult `json:"data"`
}

type CreateDraftDataResult struct {
	MongoID string `json:"mongo_id"`
	Status  string `json:"status"`
}

type SubmitAchievementResponse struct {
	Status  string                    `json:"status"`
	Message string                    `json:"message"`
	Data    SubmitAchievementData     `json:"data"`
}

type SubmitAchievementData struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type AchievementStatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"data"`
}

type AchievementDetailResponse struct {
    Status string `json:"status"`
    Data   struct {
        ID                string      `json:"id"`
        Status            string      `json:"status"`
        CreatedAt         string      `json:"created_at"`
        AchievementDetail interface{} `json:"achievement_detail"`
    } `json:"data"`
}

type UpdateAchievementResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
    Data    struct {
        ID      string `json:"id"`
        MongoID string `json:"mongo_id"`
    } `json:"data"`
}

type DeleteAchievementResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
    Data    struct {
        ID string `json:"id"`
    } `json:"data"`
}

type UploadAttachmentResponse struct {
    Status  string     `json:"status"`
    Message string     `json:"message"`
    File    Attachment `json:"file"` // pakai struct yang sama
}

type UpdateAchievementRequest struct {
	Title       *string                `json:"title"`
	Description *string                `json:"description"`
	Details     map[string]interface{} `json:"details"`
	Points      *int                   `json:"points" example:"100"`
	Tags        *[]string              `json:"tags"`
}

type RejectionNote struct {
	Note string `json:"note"`
}

type HistoryItem struct {
    Status    string    `json:"status"`
    UpdatedAt time.Time `json:"updated_at"`
    UpdatedBy string    `json:"updated_by"`
    Note      string    `json:"note,omitempty"`
}

type AchievementHistoryResponse struct {
    Status  string        `json:"status"`
    History []HistoryItem `json:"history"`
}
