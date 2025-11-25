package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	// Ambil DSN dari environment variable
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("DB_DSN tidak ditemukan, pastikan sudah diset di environment")
	}

	// Buat koneksi ke PostgreSQL
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal open koneksi: %v", err)
	}

	// Cek koneksi
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("gagal ping database: %v", err)
	}

	log.Println("Database PostgreSQL berhasil terkoneksi")
	return db, nil
}
