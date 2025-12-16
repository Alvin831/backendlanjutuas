package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	RecipientID string            `json:"recipient_id" bson:"recipient_id"` // Dosen wali ID
	SenderID    string            `json:"sender_id" bson:"sender_id"`       // Mahasiswa ID
	Type        string            `json:"type" bson:"type"`                 // achievement_submitted, achievement_verified, etc
	Title       string            `json:"title" bson:"title"`
	Message     string            `json:"message" bson:"message"`
	Data        interface{}       `json:"data" bson:"data"`                 // Additional data (achievement_id, etc)
	IsRead      bool              `json:"is_read" bson:"is_read"`
	CreatedAt   time.Time         `json:"created_at" bson:"created_at"`
	ReadAt      *time.Time        `json:"read_at,omitempty" bson:"read_at,omitempty"`
}

type CreateNotificationRequest struct {
	RecipientID string      `json:"recipient_id"`
	Type        string      `json:"type"`
	Title       string      `json:"title"`
	Message     string      `json:"message"`
	Data        interface{} `json:"data,omitempty"`
}