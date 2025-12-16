package route

import (
	"uas_backend/app/service"
	"uas_backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterNotificationRoutes(app fiber.Router) {
	notifications := app.Group("/v1/notifications")

	// GET /api/v1/notifications - Get user notifications
	notifications.Get("/", 
		middleware.AuthRequired,
		service.GetUserNotifications)

	// PUT /api/v1/notifications/:id/read - Mark notification as read
	notifications.Put("/:id/read", 
		middleware.AuthRequired,
		service.MarkNotificationAsRead)
}