package model

import "time"

type Achievement struct {
	ID              string                 `json:"id" bson:"_id,omitempty"`
	StudentID       string                 `json:"student_id" bson:"studentId"`
	AchievementType string                 `json:"achievement_type" bson:"achievementType"`
	Title           string                 `json:"title" bson:"title"`
	Description     string                 `json:"description" bson:"description"`
	Details         map[string]interface{} `json:"details" bson:"details"`
	Attachments     []Attachment           `json:"attachments" bson:"attachments"`
	Tags            []string               `json:"tags" bson:"tags"`
	Points          int                    `json:"points" bson:"points"`
	CreatedAt       time.Time              `json:"created_at" bson:"createdAt"`
	UpdatedAt       time.Time              `json:"updated_at" bson:"updatedAt"`
}

type Attachment struct {
	FileName   string    `json:"file_name" bson:"fileName"`
	FileURL    string    `json:"file_url" bson:"fileUrl"`
	FileType   string    `json:"file_type" bson:"fileType"`
	UploadedAt time.Time `json:"uploaded_at" bson:"uploadedAt"`
}

