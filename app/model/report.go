package model

import "time"

// FR-011: Achievement Statistics
type AchievementStatistics struct {
	// Basic counts
	TotalAchievements int `json:"total_achievements"`
	VerifiedCount     int `json:"verified_count"`
	PendingCount      int `json:"pending_count"`
	RejectedCount     int `json:"rejected_count"`
	DraftCount        int `json:"draft_count"`
	
	// Total prestasi per tipe/kategori
	ByCategory map[string]int `json:"by_category"`
	
	// Total prestasi per periode (bulan)
	ByPeriod map[string]int `json:"by_period"`
	
	// Top mahasiswa berprestasi
	TopStudents []TopStudent `json:"top_students"`
	
	// Distribusi tingkat kompetisi
	ByCompetitionLevel map[string]int `json:"by_competition_level"`
	
	// Summary stats
	TotalPoints       int     `json:"total_points"`
	AveragePoints     float64 `json:"average_points"`
	
	// Period info
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
}

type TopStudent struct {
	StudentID         string `json:"student_id"`
	StudentName       string `json:"student_name"`
	NIM               string `json:"nim"`
	TotalAchievements int    `json:"total_achievements"`
	TotalPoints       int    `json:"total_points"`
	VerifiedCount     int    `json:"verified_count"`
}

// Legacy structures (keep for backward compatibility)
type Statistics struct {
	TotalStudents     int `json:"total_students"`
	TotalAchievements int `json:"total_achievements"`
	TotalPoints       int `json:"total_points"`
	VerifiedCount     int `json:"verified_count"`
	PendingCount      int `json:"pending_count"`
	RejectedCount     int `json:"rejected_count"`
}

type StudentReport struct {
	StudentID         string         `json:"student_id"`
	StudentName       string         `json:"student_name"`
	NIM               string         `json:"nim"`
	Program           string         `json:"program"`
	Semester          int            `json:"semester"`
	TotalAchievements int            `json:"total_achievements"`
	TotalPoints       int            `json:"total_points"`
	AveragePoints     float64        `json:"average_points"`
	VerifiedCount     int            `json:"verified_count"`
	PendingCount      int            `json:"pending_count"`
	RejectedCount     int            `json:"rejected_count"`
	DraftCount        int            `json:"draft_count"`
	ByCategory        map[string]int `json:"by_category"`
}