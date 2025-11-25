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

	// Load file .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file tidak ditemukan, menggunakan environment system.")
	}

	// Connect database PostgreSQL
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Gagal koneksi database: ", err)
	}
	defer db.Close()

	// Init Fiber
	app := fiber.New()

	// === INIT REPO ===
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)


	// === SET REPO KE SERVICE ===
	service.SetUserRepo(userRepo)
	service.SetRoleRepo(roleRepo)
	

	// === REGISTER ROUTES ===
	api := app.Group("/api")

	route.RegisterUserRoutes(api)
	route.RegisterRoleRoutes(api)

	// Jalankan server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Server berjalan di port:", port)
	log.Fatal(app.Listen(":" + port))
}
