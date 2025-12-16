package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Achievement struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	StudentID   string            `json:"student_id" bson:"student_id"`
	Title       string            `json:"title" bson:"title"`
	Description string            `json:"description" bson:"description"`
	Category    string            `json:"category" bson:"category"`
	Points      int               `json:"points" bson:"points"`
	Status      string            `json:"status" bson:"status"` // draft, submitted, verified, rejected, deleted
	Documents   []AchievementDocument `json:"documents" bson:"documents"`
	SubmittedAt *time.Time        `json:"submitted_at,omitempty" bson:"submitted_at,omitempty"`
	VerifiedAt  *time.Time        `json:"verified_at,omitempty" bson:"verified_at,omitempty"`
	VerifiedBy  *string           `json:"verified_by,omitempty" bson:"verified_by,omitempty"`
	RejectedAt  *time.Time        `json:"rejected_at,omitempty" bson:"rejected_at,omitempty"`
	RejectionReason *string       `json:"rejection_reason,omitempty" bson:"rejection_reason,omitempty"`
	IsDeleted   bool              `json:"is_deleted" bson:"is_deleted"`
	DeletedAt   *time.Time        `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
	DeletedBy   *string           `json:"deleted_by,omitempty" bson:"deleted_by,omitempty"`
	CreatedAt   time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" bson:"updated_at"`
}

type AchievementDocument struct {
	ID       string `json:"id" bson:"id"`
	FileName string `json:"file_name" bson:"file_name"`
	FilePath string `json:"file_path" bson:"file_path"`
	FileSize int64  `json:"file_size" bson:"file_size"`
	MimeType string `json:"mime_type" bson:"mime_type"`
}

type CreateAchievementRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Category    string `json:"category" validate:"required"`
	Points      int    `json:"points" validate:"required,min=1"`
}

type UpdateAchievementRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Points      int    `json:"points"`
}

type SubmitAchievementRequest struct {
	Notes string `json:"notes,omitempty"`
}

type VerifyAchievementRequest struct {
	Notes string `json:"notes,omitempty"`
}

type RejectAchievementRequest struct {
	Reason string `json:"reason" validate:"required"`
}

type AchievementHistory struct {
	ID            string    `json:"id"`
	AchievementID string    `json:"achievement_id"`
	Status        string    `json:"status"`
	Notes         string    `json:"notes,omitempty"`
	CreatedBy     string    `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
}

type AchievementAttachment struct {
	ID            string    `json:"id"`
	AchievementID string    `json:"achievement_id"`
	FileName      string    `json:"file_name"`
	FilePath      string    `json:"file_path"`
	FileSize      int64     `json:"file_size"`
	MimeType      string    `json:"mime_type"`
	UploadedBy    string    `json:"uploaded_by"`
	CreatedAt     time.Time `json:"created_at"`
}