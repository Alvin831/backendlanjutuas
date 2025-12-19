package service

import (
	"uas_backend/app/model"
	"uas_backend/app/utils"

	"github.com/gofiber/fiber/v2"
)

// GetAllUsers godoc
// @Summary Get All Users
// @Description Get list of all users (Admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response "Users retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /v1/users [get]
func GetAllUsers(c *fiber.Ctx) error {
	users, err := userRepo.FindAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": users})
}

// GetUserByID godoc
// @Summary Get User by ID
// @Description Get user details by ID (Admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} utils.Response "User retrieved successfully"
// @Failure 404 {object} utils.Response "User not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /v1/users/{id} [get]
func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := userRepo.FindByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if user == nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "User tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"success": true, "data": user})
}

func CreateUser(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
		FullName string `json:"full_name" validate:"required"`
		RoleID   string `json:"role_id" validate:"required"`
	}
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Request tidak valid"})
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to hash password"})
	}

	// Create user object
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		RoleID:       req.RoleID,
		IsActive:     true,
	}

	newUser, err := userRepo.Create(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	// Don't return password hash in response
	newUser.PasswordHash = ""

	return c.Status(201).JSON(fiber.Map{"success": true, "data": newUser, "message": "User berhasil dibuat"})
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.User

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Request tidak valid"})
	}

	req.ID = id
	updatedUser, err := userRepo.Update(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if updatedUser == nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "User tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"success": true, "data": updatedUser, "message": "User berhasil diperbarui"})
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	deleted, err := userRepo.Delete(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if !deleted {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "User tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "User berhasil dihapus"})
}

// ================================================= UPDATE USER ROLE
func UpdateUserRole(c *fiber.Ctx) error {
	id := c.Params("id")
	
	var req struct {
		RoleID string `json:"role_id"`
	}
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Request tidak valid"})
	}

	// Get current user
	user, err := userRepo.FindByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if user == nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "User tidak ditemukan"})
	}

	// Update role
	user.RoleID = req.RoleID
	updatedUser, err := userRepo.Update(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "data": updatedUser, "message": "Role user berhasil diperbarui"})
}

// ================================================= GET USER PERMISSIONS (untuk debugging)
func GetUserPermissions(c *fiber.Ctx) error {
	userID := c.Params("userId")
	
	permissions, err := userRepo.GetUserPermissions(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user_id": userID,
			"permissions": permissions,
		},
	})
}
