package model

type Statistics struct {
	TotalStudents     int `json:"total_students"`
	TotalAchievements int `json:"total_achievements"`
	TotalPoints       int `json:"total_points"`
	VerifiedCount     int `json:"verified_count"`
	PendingCount      int `json:"pending_count"`
	RejectedCount     int `json:"rejected_count"`
}

type StudentReport struct {
	StudentID         string `json:"student_id"`
	StudentName       string `json:"student_name"`
	NIM               string `json:"nim"`
	TotalAchievements int    `json:"total_achievements"`
	TotalPoints       int    `json:"total_points"`
	VerifiedCount     int    `json:"verified_count"`
	PendingCount      int    `json:"pending_count"`
}