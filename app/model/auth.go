package model

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ==========================
// MODEL USER (untuk login)
// ==========================

// ==========================
// MODEL UNTUK LOGIN REQUEST
// ==========================
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ==========================
// LOGIN RESPONSE
// Sesuai SRS: access_token + refresh_token
// ==========================
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ==========================
// MODEL USER DENGAN PERMISSIONS
// untuk JWT payload
// ==========================
type UserWithPermissions struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions"`
}

// ==========================
// JWT CLAIMS
// Sudah cocok untuk JWT di utils/jwt.go
// ==========================
type JWTClaims struct {
	UserID      uuid.UUID `json:"user_id"`
	Username    string    `json:"username"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions"`
	jwt.RegisteredClaims
}

// ==========================
// REFRESH TOKEN REQUEST
// untuk endpoint refresh
// ==========================
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// ==========================
// PROFILE RESPONSE
// digunakan ketika user melihat profilnya
// ==========================
type ProfileResponse struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions"`
}
