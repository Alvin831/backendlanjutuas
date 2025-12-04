package service

import (
	"uas_backend/app/model"

	"github.com/gofiber/fiber/v2"
)

func GetAllUsers(c *fiber.Ctx) error {
	users, err := userRepo.FindAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": users})
}

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
	var req model.User
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Request tidak valid"})
	}

	newUser, err := userRepo.Create(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

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
