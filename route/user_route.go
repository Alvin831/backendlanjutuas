package route

import (
	"time"
	"uas_backend/app/service"
	"uas_backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app fiber.Router) {
	user := app.Group("/v1/users")

	// Testing endpoint tanpa permission check (temporary)
	user.Get("/test", 
		middleware.AuthRequired, 
		service.GetAllUsers)
		
	// Debug: Check user permissions (temporary)
	user.Get("/permissions/:userId", 
		middleware.AuthRequired, 
		service.GetUserPermissions)

	// 5.2 Users (Admin) - Sesuai dengan endpoint specification
	user.Get("/", 
		middleware.AuthRequired, 
		middleware.RateLimitByUser(100, time.Hour),
		middleware.PermissionRequired("manage_users"), 
		service.GetAllUsers)
		
	user.Get("/:id", 
		middleware.AuthRequired, 
		middleware.PermissionRequired("manage_users"), 
		service.GetUserByID)
		
	user.Post("/", 
		middleware.AuthRequired, 
		middleware.RateLimitByPermission("manage_users", 10, time.Minute),
		middleware.PermissionRequired("manage_users"), 
		service.CreateUser)
		
	user.Put("/:id", 
		middleware.AuthRequired, 
		middleware.PermissionRequired("manage_users"), 
		service.UpdateUser)
		
	user.Delete("/:id", 
		middleware.AuthRequired, 
		middleware.RateLimitByPermission("manage_users", 5, time.Minute),
		middleware.PermissionRequired("manage_users"), 
		service.DeleteUser)

	// PUT /api/v1/users/:id/role - Update user role
	user.Put("/:id/role", 
		middleware.AuthRequired, 
		middleware.PermissionRequired("manage_users"), 
		service.UpdateUserRole)
}
