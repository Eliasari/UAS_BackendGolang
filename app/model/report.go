package model

// ==========================
// ITEM GENERIC
// ==========================
type StatisticItem struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Total int    `json:"total"`
}

// ==========================
// TOP STUDENT ITEM
// ==========================
type StatisticTopStudentItem struct {
	StudentID   string `json:"student_id"`
	StudentName string `json:"student_name"`
	Total       int    `json:"total"`
}

// ==========================
// RESPONSE UTAMA
// ==========================
type StatisticsResponse struct {
	TotalPerType     []StatisticItem           `json:"total_per_type"`
	TotalPerPeriod   []StatisticItem           `json:"total_per_period"`
	CompetitionLevel []StatisticItem           `json:"competition_levels"`
	TopStudents      []StatisticTopStudentItem `json:"top_students"`
}
