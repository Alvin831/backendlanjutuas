package model

import "time"

// Achievement reference in PostgreSQL for reporting and analytics
type AchievementReference struct {
	ID            string    `json:"id"`
	MongoID       string    `json:"mongo_id"`       // Reference to MongoDB ObjectID
	StudentID     string    `json:"student_id"`
	Title         string    `json:"title"`
	Category      string    `json:"category"`
	Points        int       `json:"points"`
	Status        string    `json:"status"`
	IsDeleted     bool      `json:"is_deleted"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}