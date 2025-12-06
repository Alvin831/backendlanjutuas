package main

import (
	"log"
	"os"

	"uas_backend/app/repository"
	"uas_backend/app/service"
	"uas_backend/database"
	"uas_backend/middleware"
	"uas_backend/route"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file tidak ditemukan, menggunakan environment system.")
	}

	// Connect DB
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Gagal koneksi database: ", err)
	}
	defer db.Close()

	// Init Fiber
	app := fiber.New()

	// ========= INIT REPOSITORY =========
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)

	// ========= SET REPO TO SERVICE =========
	service.SetUserRepo(userRepo)
	service.SetRoleRepo(roleRepo)
	service.SetPermissionRepo(permissionRepo)

	// ========= INIT MIDDLEWARE =========
	authMw := middleware.NewAuthMiddleware(roleRepo)

	// ========= REGISTER ROUTES =========
	api := app.Group("/api")

	// 🔹 Route Authentication (tidak pakai middleware)
	route.RegisterAuthRoutes(api)

	// 🔹 Route lain harus login dulu
	protected := api.Use(middleware.AuthRequired)

	// User routes (butuh token)
	route.RegisterUserRoutes(protected, authMw)

	// Role routes (butuh token)
	route.RegisterRoleRoutes(protected, authMw)

	// Permission routes (butuh token)
	route.RegisterPermissionRoutes(protected)

	// Run server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Server berjalan di port:", port)
	log.Fatal(app.Listen(":" + port))
}
