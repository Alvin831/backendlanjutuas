package route

import (
	"time"
	"uas_backend/app/service"
	"uas_backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterReportRoutes(app fiber.Router) {
	reports := app.Group("/v1/reports")

	// 5.8 Reports & Analytics - Sesuai dengan endpoint specification
	
	// GET /api/v1/reports/statistics
	reports.Get("/statistics", 
		middleware.AuthRequired,
		middleware.RateLimitByUser(50, time.Hour),
		middleware.AnyPermissionRequired("view_all", "verify_prestasi", "manage_users"),
		service.GetStatistics)

	// GET /api/v1/reports/student/:id
	reports.Get("/student/:id", 
		middleware.AuthRequired,
		middleware.AnyPermissionRequired("view_all", "verify_prestasi", "manage_users"),
		service.GetStudentReport)
}