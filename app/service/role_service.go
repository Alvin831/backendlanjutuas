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
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    roles,
		"message": "Data role berhasil diambil",
	})
}

func GetRoleByID(c *fiber.Ctx) error {
	id := c.Params("id")
	role, err := roleRepo.FindByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if role == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Role tidak ditemukan"})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    role,
		"message": "Data role berhasil diambil",
	})
}

func CreateRole(c *fiber.Ctx) error {
	var req model.Role
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request tidak valid"})
	}
	role, err := roleRepo.Create(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"data":    role,
		"message": "Role berhasil dibuat",
	})
}
