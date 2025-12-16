package model

type CreateAchievementRequest struct {
	AchievementType string                 `json:"achievement_type"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Details         map[string]interface{} `json:"details"`
	Points          int                    `json:"points" example:"50"`
	Tags            []string               `json:"tags"`
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
