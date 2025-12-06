package route

import (
	"uas_backend/app/service"
	"uas_backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	user := app.Group("/users")

	user.Get("/", middleware.AuthRequired, authMiddleware.PermissionRequired("manage_users"), service.GetAllUsers)
	user.Get("/:id", middleware.AuthRequired, authMiddleware.PermissionRequired("manage_users"), service.GetUserByID)
	user.Post("/", middleware.AuthRequired, authMiddleware.PermissionRequired("manage_users"), service.CreateUser)
	user.Put("/:id", middleware.AuthRequired, authMiddleware.PermissionRequired("manage_users"), service.UpdateUser)
	user.Delete("/:id", middleware.AuthRequired, authMiddleware.PermissionRequired("manage_users"), service.DeleteUser)
}
