package route

import (
	"uas_backend/app/service"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoleRoutes(app fiber.Router) {
	app.Get("/roles", service.GetAllRoles)
	app.Get("/roles/:id", service.GetRoleByID)
	app.Post("/roles", service.CreateRole)
}
