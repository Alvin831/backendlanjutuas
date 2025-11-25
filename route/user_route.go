package route

import (
	"uas_backend/app/service"
	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app fiber.Router) {
	app.Get("/users", service.GetAllUsers)
	app.Get("/users/:id", service.GetUserByID)
	app.Post("/users", service.CreateUser)
	app.Put("/users/:id", service.UpdateUser)
	app.Delete("/users/:id", service.DeleteUser)
}
