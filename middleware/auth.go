package middleware

import (
	"strings"
	"uas_backend/app/utils"

	"github.com/gofiber/fiber/v2"
)

// ==============================
// JWT VALIDATION (WAJIB LOGIN)
// ==============================
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token akses diperlukan",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Format token tidak valid (gunakan: Bearer <token>)",
			})
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token tidak valid atau sudah expired",
			})
		}

		// Simpan claims ke context agar bisa digunakan oleh controller lain
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)
		c.Locals("permissions", claims.Permissions)

		return c.Next()
	}
}

// ==============================
// ROLE-BASED AUTH (SRS: RBAC)
// contoh: RoleRequired("admin")
// ==============================
func RoleRequired(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		role, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "User tidak valid"})
		}

		for _, r := range allowedRoles {
			if r == role {
				return c.Next() // role cocok → lanjut
			}
		}

		return c.Status(403).JSON(fiber.Map{
			"error": "Akses ditolak. Role tidak memiliki izin",
		})
	}
}

// ==============================
// PERMISSION-BASED AUTH
// sesuai SRS (FR-002 Access Control)
// example: PermissionRequired("achievement.create")
// ==============================
func PermissionRequired(requiredPermission string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		permissions, ok := c.Locals("permissions").([]string)
		if !ok {
			return c.Status(403).JSON(fiber.Map{
				"error": "Tidak dapat membaca permissions",
			})
		}

		for _, p := range permissions {
			if p == requiredPermission {
				return c.Next() // permission cocok → lanjut
			}
		}

		return c.Status(403).JSON(fiber.Map{
			"error": "Akses ditolak. Permission tidak mencukupi",
		})
	}
}
