package route

import (
	"uas_backend/app/service"
	"uas_backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(api fiber.Router) {
	auth := api.Group("/auth")

	// Public (no token)
	auth.Post("/register", service.Register)
	auth.Post("/login", service.Login)

	// Protected (token required)
	auth.Get("/profile", middleware.AuthRequired, service.GetProfile)
}
