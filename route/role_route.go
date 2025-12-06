package route

import (
	"uas_backend/app/service"
	"uas_backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoleRoutes(app fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	role := app.Group("/roles")

	role.Get("/", middleware.AuthRequired, authMiddleware.PermissionRequired("manage_roles"), service.GetAllRoles)
	role.Get("/:id", middleware.AuthRequired, authMiddleware.PermissionRequired("manage_roles"), service.GetRoleByID)
	role.Post("/", middleware.AuthRequired, authMiddleware.PermissionRequired("manage_roles"), service.CreateRole)

    // Mapping permission ke role
	role.Post("/:roleId/permissions", middleware.AuthRequired, authMiddleware.PermissionRequired("manage_roles"), service.AssignPermissionsToRole)
}
