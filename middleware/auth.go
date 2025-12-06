package middleware

import (
	"strings"
	"uas_backend/app/model"
	"uas_backend/app/repository"
	"uas_backend/app/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	roleRepo *repository.RoleRepository
}

func NewAuthMiddleware(roleRepo *repository.RoleRepository) *AuthMiddleware {
	return &AuthMiddleware{roleRepo: roleRepo}
}

// ====================================================================
// 1️⃣ AUTH REQUIRED — Validasi TOKEN dulu
// ====================================================================
func AuthRequired(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(model.WebResponse{
			Code:    401,
			Status:  "error",
			Message: "Missing authorization header",
		})
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(401).JSON(model.WebResponse{
			Code:    401,
			Status:  "error",
			Message: "Invalid token format",
		})
	}

	claims, err := utils.ParseToken(tokenParts[1])
	if err != nil {
		return c.Status(401).JSON(model.WebResponse{
			Code:    401,
			Status:  "error",
			Message: "Invalid or expired token",
		})
	}

	// Simpan ke context untuk dipakai service berikutnya
	c.Locals("user_id", claims.UserID)
	c.Locals("role", claims.Role)
	c.Locals("permissions", claims.Permissions)

	return c.Next()
}

// ====================================================================
// 2️⃣ PERMISSION REQUIRED — Cek apakah user punya izin yang diperlukan
// ====================================================================
func (m *AuthMiddleware) PermissionRequired(requiredPerm string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		userPermsInterface := c.Locals("permissions")
		if userPermsInterface == nil {
			return c.Status(403).JSON(model.WebResponse{
				Code:    403,
				Status:  "error",
				Message: "No permissions found",
			})
		}

		// Convert interface agar menjadi []string
		var userPerms []string
		switch v := userPermsInterface.(type) {
		case []string:
			userPerms = v
		case []interface{}:
			for _, item := range v {
				if s, ok := item.(string); ok {
					userPerms = append(userPerms, s)
				}
			}
		}

		// Check match
		allowed := false
		for _, p := range userPerms {
			if p == requiredPerm {
				allowed = true
				break
			}
		}

		if !allowed {
			return c.Status(403).JSON(model.WebResponse{
				Code:    403,
				Status:  "error",
				Message: "Access denied — missing permission: " + requiredPerm,
			})
		}

		return c.Next()
	}
}
