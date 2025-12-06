package service

import (
	"uas_backend/app/model"
	"uas_backend/app/repository"
	"github.com/gofiber/fiber/v2"
)

var permRepo *repository.PermissionRepository

// Dipanggil dari main.go
func SetPermissionRepo(repo *repository.PermissionRepository) {
	permRepo = repo
}

// ===================== GET ALL PERMISSIONS =====================
func GetAllPermissions(c *fiber.Ctx) error {
	permissions, err := permRepo.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    permissions,
	})
}

// ===================== GET PERMISSION BY ID =====================
func GetPermissionByID(c *fiber.Ctx) error {
	id := c.Params("id")

	permission, err := permRepo.FindByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if permission == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Permission tidak ditemukan"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    permission,
	})
}

// ===================== CREATE PERMISSION =====================
func CreatePermission(c *fiber.Ctx) error {
	var req model.Permission
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Bad request"})
	}

	created, err := permRepo.Create(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Permission berhasil dibuat",
		"data":    created,
	})
}

// ===================== DELETE PERMISSION =====================
func DeletePermission(c *fiber.Ctx) error {
	id := c.Params("id")

	deleted, err := permRepo.Delete(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if !deleted {
		return c.Status(404).JSON(fiber.Map{"error": "Permission tidak ditemukan"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Permission berhasil dihapus",
	})
}

// ===================== ASSIGN PERMISSION TO ROLE =====================
func AssignPermission(c *fiber.Ctx) error {
	roleID := c.Params("role_id")
	permID := c.Params("permission_id")

	err := permRepo.Assign(roleID, permID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Permission berhasil ditambahkan ke role",
	})
}

// ===================== REMOVE PERMISSION FROM ROLE =====================
func RemovePermission(c *fiber.Ctx) error {
	roleID := c.Params("role_id")
	permID := c.Params("permission_id")

	err := permRepo.Remove(roleID, permID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Permission berhasil dihapus dari role",
	})
}
