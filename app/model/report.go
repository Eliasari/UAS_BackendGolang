package model

type StatisticItem struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Total int    `json:"total"`
}

type StatisticTopStudentItem struct {
	StudentID   string `json:"student_id"`
	StudentName string `json:"student_name"`
	Total       int    `json:"total"`
}

type StatisticsResponse struct {
	TotalPerType     []StatisticItem           `json:"total_per_type"`
	TotalPerPeriod   []StatisticItem           `json:"total_per_period"`
	CompetitionLevel []StatisticItem           `json:"competition_levels"`
	TopStudents      []StatisticTopStudentItem `json:"top_students"`
}

type StudentAchievementReportResponse struct {
	ID              interface{}            `json:"_id"`
	StudentID       string                 `json:"studentId"`
	AchievementType string                 `json:"achievementType"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Details         map[string]interface{} `json:"details"`
	Points          int                    `json:"points"`
	CreatedAt       interface{}            `json:"createdAt"`
}
