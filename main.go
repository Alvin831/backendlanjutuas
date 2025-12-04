package main

import (
	"log"
	"os"

	"uas_backend/database"
	"uas_backend/app/repository"
	"uas_backend/app/service"
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

	// Init repository
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)

	// Set repo to service
	service.SetUserRepo(userRepo)
	service.SetRoleRepo(roleRepo)

	// Register routes
	api := app.Group("/api")
	route.RegisterAuthRoutes(api)  // ⬅ WAJIB agar login & register aktif
	route.RegisterUserRoutes(api)
	route.RegisterRoleRoutes(api)

	// Run server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Server berjalan di port:", port)
	log.Fatal(app.Listen(":" + port))
}
