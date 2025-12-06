package route

import (
	"uas_backend/app/service"
	"github.com/gofiber/fiber/v2"
)

func RegisterPermissionRoutes(api fiber.Router) {

	// Permission CRUD
	api.Get("/permissions", service.GetAllPermissions)
	api.Get("/permissions/:id", service.GetPermissionByID)
	api.Post("/permissions", service.CreatePermission)
	api.Delete("/permissions/:id", service.DeletePermission)

	// Assign & Remove permission to role
	api.Post("/roles/:role_id/permissions/:permission_id", service.AssignPermission)
	api.Delete("/roles/:role_id/permissions/:permission_id", service.RemovePermission)
}
