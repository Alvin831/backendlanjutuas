package service

import (
	"uas_backend/app/model"
	"uas_backend/app/repository"
	"github.com/gofiber/fiber/v2"
)

var userRepo *repository.UserRepository

func SetUserRepo(repo *repository.UserRepository) {
	userRepo = repo
}

func GetAllUsers(c *fiber.Ctx) error {
	users, err := userRepo.FindAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    users,
		"message": "Data user berhasil diambil",
	})
}

func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := userRepo.FindByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if user == nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    user,
		"message": "Data user berhasil diambil",
	})
}

func CreateUser(c *fiber.Ctx) error {
	var req model.User
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request body tidak valid"})
	}
	newUser, err := userRepo.Create(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"data":    newUser,
		"message": "User berhasil dibuat",
	})
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.User
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request tidak valid"})
	}
	req.ID = id

	updatedUser, err := userRepo.Update(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if updatedUser == nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    updatedUser,
		"message": "User berhasil diperbarui",
	})
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	deleted, err := userRepo.Delete(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if !deleted {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User berhasil dihapus",
	})
}
