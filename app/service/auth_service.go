package service

import (
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

// ================================================= REGISTER
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

// ================================================= LOGIN
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

	token, _ := utils.GenerateToken(user.ID, user.RoleID, []string{})

	return c.JSON(utils.SuccessResponse(
		"Login berhasil", 200, fiber.Map{
			"access_token": token,
		},
	))
}

// ================================================= PROFILE
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
