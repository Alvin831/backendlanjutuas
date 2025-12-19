package model

import "time"

type Student struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	NIM       string    `json:"nim"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Program   string    `json:"program"`
	Semester  int       `json:"semester"`
	AdvisorID *string   `json:"advisor_id,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Lecturer struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	LecturerID string    `json:"lecturer_id"` // NIDN or employee ID
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Department string    `json:"department"`  // Faculty/Department
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UpdateAdvisorRequest struct {
	AdvisorID string `json:"advisor_id" validate:"required"`
}