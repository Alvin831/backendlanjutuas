package route

import (
	"uas_backend/app/service"
	"uas_backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(api fiber.Router) {
	auth := api.Group("/v1/auth")

	// 5.1 Authentication - Sesuai dengan endpoint specification
	auth.Post("/login", service.Login)
	auth.Post("/refresh", service.RefreshToken)
	auth.Post("/logout", middleware.AuthRequired, service.Logout)
	auth.Get("/profile", middleware.AuthRequired, service.GetProfile)
}
