package service

import (
	"uas_backend/app/model"
	"uas_backend/app/repository"

	"github.com/gofiber/fiber/v2"
)

var roleRepo *repository.RoleRepository

func SetRoleRepo(repo *repository.RoleRepository) {
	roleRepo = repo
}

func GetAllRoles(c *fiber.Ctx) error {
	roles, err := roleRepo.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": roles})
}

func GetRoleByID(c *fiber.Ctx) error {
	id := c.Params("id")
	role, err := roleRepo.FindByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if role == nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Role tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"success": true, "data": role})
}

func CreateRole(c *fiber.Ctx) error {
	var req model.Role
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Request tidak valid"})
	}

	role, err := roleRepo.Create(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "data": role, "message": "Role berhasil dibuat"})
}

// ================== ASSIGN PERMISSION TO ROLE ==================
func AssignPermissionsToRole(c *fiber.Ctx) error {
	roleId := c.Params("roleId")
	if roleId == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "roleId diperlukan"})
	}

	var req struct {
		PermissionIDs []string `json:"permission_ids"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Request body tidak valid"})
	}

	if len(req.PermissionIDs) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "permission_ids tidak boleh kosong"})
	}

	err := roleRepo.AssignPermissions(roleId, req.PermissionIDs)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Permissions berhasil ditambahkan ke role",
		"data": fiber.Map{
			"role_id":        roleId,
			"permission_ids": req.PermissionIDs,
		},
	})
}

