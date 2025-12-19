package model

import (
	"github.com/golang-jwt/jwt/v5"
)

// LoginRequest represents login credentials
// @Description Login request payload
type LoginRequest struct {
	Username string `json:"username" example:"admin" binding:"required"`
	Password string `json:"password" example:"123456" binding:"required"`
}

// RegisterRequest represents registration data
// @Description User registration request payload
type RegisterRequest struct {
	Username string `json:"username" example:"newuser" binding:"required"`
	Email    string `json:"email" example:"user@example.com" binding:"required"`
	FullName string `json:"full_name" example:"New User" binding:"required"`
	Password string `json:"password" example:"123456" binding:"required"`
	RoleID   string `json:"role_id" example:"f464ceb1-5481-49cf-99f0-d8f2d66f4506" binding:"required"`
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
