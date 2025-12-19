package main

import (
	"log"
	"os"
	"time"

	"uas_backend/app/repository"
	"uas_backend/app/service"
	"uas_backend/database"
	"uas_backend/middleware"
	"uas_backend/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	_ "uas_backend/docs" // Import generated docs
)

// @title Achievement Management API
// @version 1.0
// @description API untuk sistem manajemen prestasi mahasiswa dengan RBAC
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {

	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file tidak ditemukan, menggunakan environment system.")
	}

	// Connect PostgreSQL
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Gagal koneksi PostgreSQL: ", err)
	}
	defer db.Close()

	// Connect MongoDB
	err = database.ConnectMongoDB()
	if err != nil {
		log.Fatal("Gagal koneksi MongoDB: ", err)
	}

	// Init Fiber
	app := fiber.New()

	// ========= GLOBAL MIDDLEWARE =========
	// Audit logging untuk semua requests (sudah include rate limiting)
	app.Use(middleware.AuditMiddleware())

	// ========= INIT REPOSITORY =========
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	achievementRepo := repository.NewAchievementRepository()
	achievementRefRepo := repository.NewAchievementReferenceRepository(db)
	studentRepo := repository.NewStudentRepository(db)

	// ========= SET REPO TO SERVICE =========
	service.SetUserRepo(userRepo)
	service.SetRoleRepo(roleRepo)
	service.SetAchievementRepo(achievementRepo)
	service.SetAchievementReferenceRepo(achievementRefRepo)
	service.SetStudentRepo(studentRepo)
	

	// ========= SWAGGER DOCUMENTATION =========
	app.Get("/swagger/*", swagger.HandlerDefault) // default
	app.Get("/docs/*", swagger.New(swagger.Config{ // custom
		URL:         "http://localhost:3000/swagger/doc.json",
		DeepLinking: false,
		DocExpansion: "none",
	}))

	// ========= REGISTER ROUTES =========
	api := app.Group("/api")

	// 🔹 Auth routes (LOGIN, REGISTER) → TANPA middleware
	route.RegisterAuthRoutes(api)

	// 🔹 Protected routes → WAJIB TOKEN
	protected := api.Group("", middleware.AuthRequired)

	// 🔹 User routes
	route.RegisterUserRoutes(protected)

	// 🔹 Role routes
	route.RegisterRoleRoutes(protected)

	// 🔹 Achievement routes
	route.RegisterAchievementRoutes(protected)

	// 🔹 Student & Lecturer routes
	route.RegisterStudentRoutes(protected)

	// 🔹 Report routes
	route.RegisterReportRoutes(protected)

	// 🔹 Notification routes
	route.RegisterNotificationRoutes(protected)



	// Run server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Server berjalan di port:", port)
	log.Fatal(app.Listen(":" + port))
}
