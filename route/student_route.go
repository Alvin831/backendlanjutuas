package route

import (
	"time"
	"uas_backend/app/service"
	"uas_backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterStudentRoutes(app fiber.Router) {
	// 5.5 Students & Lecturers - Sesuai dengan endpoint specification
	
	// Students endpoints
	students := app.Group("/v1/students")
	
	// GET /api/v1/students
	students.Get("/", 
		middleware.AuthRequired,
		middleware.RateLimitByUser(100, time.Hour),
		middleware.AnyPermissionRequired("view_students", "manage_students"),
		service.GetAllStudents)

	// GET /api/v1/students/:id
	students.Get("/:id", 
		middleware.AuthRequired,
		middleware.AnyPermissionRequired("view_students", "manage_students"),
		service.GetStudentByID)

	// GET /api/v1/students/:id/achievements
	students.Get("/:id/achievements", 
		middleware.AuthRequired,
		middleware.AnyPermissionRequired("view_students", "manage_students", "view_achievements"),
		service.GetStudentAchievements)

	// PUT /api/v1/students/:id/advisor
	students.Put("/:id/advisor", 
		middleware.AuthRequired,
		middleware.PermissionRequired("manage_students"),
		service.UpdateStudentAdvisor)

	// Lecturers endpoints
	lecturers := app.Group("/v1/lecturers")
	
	// GET /api/v1/lecturers
	lecturers.Get("/", 
		middleware.AuthRequired,
		middleware.RateLimitByUser(100, time.Hour),
		middleware.AnyPermissionRequired("view_lecturers", "manage_lecturers"),
		service.GetAllLecturers)

	// GET /api/v1/lecturers/:id/advisees
	lecturers.Get("/:id/advisees", 
		middleware.AuthRequired,
		middleware.AnyPermissionRequired("view_lecturers", "manage_lecturers", "verify_achievements"),
		service.GetLecturerAdvisees)
}