package route

import (
	// Import service package Anda, misal: "GOLANG/Domain/service"
	"uas_backend/app/service"
	"uas_backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoleRoutes(app fiber.Router) { // Hapus koma berlebih di parameter
	role := app.Group("/roles")

	// PENTING: Urutannya harus Auth dulu, baru Permission

	// 1. Get All Roles
	role.Get("/", 
		middleware.AuthRequired, 
		middleware.PermissionRequired("manage_roles"), 
		service.GetAllRoles,
	)

	// 2. Get Role By ID
	role.Get("/:id", 
		middleware.AuthRequired, 
		middleware.PermissionRequired("manage_roles"), 
		service.GetRoleByID,
	)

	// 3. Create Role
	role.Post("/", 
		middleware.AuthRequired, 
		middleware.PermissionRequired("manage_roles"), 
		service.CreateRole,
	)

	// 4. Mapping permission ke role
	role.Post("/:roleId/permissions", 
		middleware.AuthRequired, 
		middleware.PermissionRequired("manage_roles"), 
		service.AssignPermissionsToRole,
	)
}