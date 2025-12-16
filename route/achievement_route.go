package route

import (
	"time"
	"uas_backend/app/service"
	"uas_backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterAchievementRoutes(app fiber.Router) {
	achievement := app.Group("/v1/achievements")

	// 5.4 Achievements - Sesuai dengan endpoint specification
	
	// GET /api/v1/achievements - List (filtered by role)
	achievement.Get("/", 
		middleware.AuthRequired,
		middleware.RateLimitByUser(200, time.Hour),
		middleware.AnyPermissionRequired("view_all", "create_prestasi", "verify_prestasi"),
		service.GetAllAchievements)

	// GET /api/v1/achievements/advisees - View Prestasi Mahasiswa Bimbingan (FR-006)
	// IMPORTANT: This must come BEFORE /:id route to avoid conflicts
	achievement.Get("/advisees", 
		middleware.AuthRequired,
		middleware.RateLimitByUser(50, time.Hour),
		middleware.PermissionRequired("verify_prestasi"), // Hanya dosen wali yang bisa akses
		service.GetAdviseesAchievements)

	// GET /api/v1/achievements/:id - Detail
	achievement.Get("/:id", 
		middleware.AuthRequired,
		middleware.AnyPermissionRequired("view_all", "create_prestasi", "verify_prestasi"),
		service.GetAchievementByID)

	// POST /api/v1/achievements - Create (Mahasiswa)
	achievement.Post("/", 
		middleware.AuthRequired,
		middleware.RateLimitByPermission("create_prestasi", 10, time.Hour),
		middleware.PermissionRequired("create_prestasi"),
		service.CreateAchievement)

	// PUT /api/v1/achievements/:id - Update (Mahasiswa)
	achievement.Put("/:id", 
		middleware.AuthRequired,
		middleware.PermissionRequired("update_prestasi"),
		service.UpdateAchievement)

	// DELETE /api/v1/achievements/:id - Delete (Mahasiswa)
	achievement.Delete("/:id", 
		middleware.AuthRequired,
		middleware.RateLimitByPermission("delete_prestasi", 5, time.Hour),
		middleware.PermissionRequired("delete_prestasi"),
		service.DeleteAchievement)

	// POST /api/v1/achievements/:id/submit - Submit for verification
	achievement.Post("/:id/submit", 
		middleware.AuthRequired,
		middleware.PermissionRequired("create_prestasi"),
		service.SubmitAchievement)

	// POST /api/v1/achievements/:id/verify - Verify (Dosen Wali)
	achievement.Post("/:id/verify", 
		middleware.AuthRequired,
		middleware.PermissionRequired("verify_prestasi"),
		service.VerifyAchievement)

	// POST /api/v1/achievements/:id/reject - Reject (Dosen Wali)
	achievement.Post("/:id/reject", 
		middleware.AuthRequired,
		middleware.PermissionRequired("verify_prestasi"),
		service.RejectAchievement)

	// GET /api/v1/achievements/:id/history - Status history
	achievement.Get("/:id/history", 
		middleware.AuthRequired,
		middleware.AnyPermissionRequired("view_all", "create_prestasi", "verify_prestasi"),
		service.GetAchievementHistory)

	// POST /api/v1/achievements/:id/attachments - Upload files
	achievement.Post("/:id/attachments", 
		middleware.AuthRequired,
		middleware.RateLimitByPermission("create_prestasi", 20, time.Hour),
		middleware.PermissionRequired("create_prestasi"),
		service.UploadAchievementAttachment)
}