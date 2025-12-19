package service

import (
	"strings"
	"uas_backend/app/model"
	"uas_backend/app/repository"
	"uas_backend/app/utils"

	"github.com/gofiber/fiber/v2"
)

var userRepo *repository.UserRepository

// dipanggil dari main.go
func SetUserRepo(repo *repository.UserRepository) {
	userRepo = repo
}

// Register godoc
// @Summary User Registration
// @Description Register a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "Registration data"
// @Success 201 {object} utils.Response "User created successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /v1/auth/register [post]
func Register(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Bad request", 400, nil))
	}

	hash, _ := utils.HashPassword(req.Password)

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		FullName:     req.FullName,
		PasswordHash: hash,
		RoleID:       req.RoleID,
		IsActive:     true,
	}

	createdUser, err := userRepo.Create(user)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse(err.Error(), 500, nil))
	}

	return c.Status(201).JSON(utils.SuccessResponse(
		"User berhasil dibuat", 201, fiber.Map{
			"id":        createdUser.ID,
			"username":  createdUser.Username,
			"full_name": createdUser.FullName,
			"email":     createdUser.Email,
			"role_id":   createdUser.RoleID,
		},
	))
}

// Login godoc
// @Summary User Login
// @Description Authenticate user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login credentials"
// @Success 200 {object} utils.Response "Login successful"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Invalid credentials"
// @Router /v1/auth/login [post]
func Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Bad request", 400, nil))
	}

	user, err := userRepo.FindByUsername(req.Username)
	if err != nil || user == nil {
		return c.Status(401).JSON(utils.ErrorResponse("Username atau password salah", 401, nil))
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return c.Status(401).JSON(utils.ErrorResponse("Username atau password salah", 401, nil))
	}

	// Load user permissions
	permissions, err := userRepo.GetUserPermissions(user.ID)
	if err != nil {
		// Jika error, tetap lanjut dengan empty permissions
		permissions = []string{}
	}

	token, _ := utils.GenerateToken(user.ID, user.RoleID, permissions)

	return c.JSON(utils.SuccessResponse(
		"Login berhasil", 200, fiber.Map{
			"access_token": token,
			"permissions":  permissions, // Include permissions in response
		},
	))
}

// GetProfile godoc
// @Summary Get User Profile
// @Description Get current user profile information
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response "Profile retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "User not found"
// @Router /v1/auth/profile [get]
func GetProfile(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(401).JSON(utils.ErrorResponse("Authorization header diperlukan", 401, nil))
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	claims, err := utils.ParseToken(token)
	if err != nil {
		return c.Status(401).JSON(utils.ErrorResponse("Token tidak valid atau expired", 401, nil))
	}

	user, err := userRepo.FindByID(claims.UserID)
	if err != nil || user == nil {
		return c.Status(404).JSON(utils.ErrorResponse("User tidak ditemukan", 404, nil))
	}

	return c.JSON(utils.SuccessResponse(
		"Profile berhasil diambil", 200, model.ProfileResponse{
			ID:          user.ID,
			Username:    user.Username,
			FullName:    user.FullName,
			Email:       user.Email,
			Role:        user.RoleID,
			Permissions: claims.Permissions,
		},
	))
}

// RefreshToken godoc
// @Summary Refresh Access Token
// @Description Refresh expired access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} utils.Response "Token refreshed successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Invalid refresh token"
// @Router /v1/auth/refresh [post]
func RefreshToken(c *fiber.Ctx) error {
	var req model.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Bad request", 400, nil))
	}

	// Parse refresh token (for now, we'll use the same JWT parsing)
	claims, err := utils.ParseToken(req.RefreshToken)
	if err != nil {
		return c.Status(401).JSON(utils.ErrorResponse("Invalid refresh token", 401, nil))
	}

	// Get user from database to ensure still active
	user, err := userRepo.FindByID(claims.UserID)
	if err != nil || user == nil || !user.IsActive {
		return c.Status(401).JSON(utils.ErrorResponse("User not found or inactive", 401, nil))
	}

	// Load fresh permissions
	permissions, err := userRepo.GetUserPermissions(user.ID)
	if err != nil {
		permissions = []string{}
	}

	// Generate new access token
	newToken, err := utils.GenerateToken(user.ID, user.RoleID, permissions)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to generate token", 500, nil))
	}

	return c.JSON(utils.SuccessResponse("Token refreshed successfully", 200, fiber.Map{
		"access_token": newToken,
		"permissions":  permissions,
	}))
}

// Logout godoc
// @Summary User Logout
// @Description Logout user and invalidate token
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response "Logout successful"
// @Router /v1/auth/logout [post]
func Logout(c *fiber.Ctx) error {
	// Get token from header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(utils.ErrorResponse("Missing authorization header", 401, nil))
	}

	// Extract token
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(401).JSON(utils.ErrorResponse("Invalid token format", 401, nil))
	}

	token := tokenParts[1]
	
	// Add token to blacklist (simple in-memory for now)
	// In production, use Redis or database
	blacklistToken(token)

	return c.JSON(utils.SuccessResponse("Logout berhasil", 200, nil))
}

// Simple in-memory token blacklist (in production, use Redis)
var tokenBlacklist = make(map[string]bool)

func blacklistToken(token string) {
	tokenBlacklist[token] = true
}

func IsTokenBlacklisted(token string) bool {
	return tokenBlacklist[token]
}
