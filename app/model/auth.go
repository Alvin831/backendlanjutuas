package model

import (
	"github.com/golang-jwt/jwt/v5"
)

// REQUEST untuk login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// REGISTER REQUEST (kalau ada register)
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
	RoleID   string `json:"role_id"`
}

// LOGIN RESPONSE — sesuai SRS
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// JWT CLAIMS — nilai yang disimpan di token
type JWTClaims struct {
	UserID      string   `json:"user_id"`
	Username    string   `json:"username"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// REQUEST untuk refresh token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RESPONSE untuk endpoint profile
type ProfileResponse struct {
	ID          string   `json:"id"`
	FullName    string   `json:"full_name"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

// WebResponse adalah struktur standar untuk response JSON (used in middleware)
type WebResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
