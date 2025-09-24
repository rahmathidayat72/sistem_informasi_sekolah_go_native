package main

import (
	"go_rest_native_sekolah/config"
	"go_rest_native_sekolah/router"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("[INFO] 🚀 Aplikasi dimulai...")

	// Load environment variables
	config.LoadEnv()

	// Ambil PORT dari environment
	port := os.Getenv("PORT")
	if port == "" {
		// PORT tidak ditemukan di environment variable
		log.Fatal("[FATAL] ❌ PORT tidak ditemukan di environment variable")
	}

	// Inisialisasi koneksi database
	db, err := config.InitPostgreSQLPool()
	if err != nil {
		log.Fatalf("[FATAL] ❌ Gagal terhubung ke database: %v", err)
	}
	defer func() {
		db.Close()
		log.Println("[INFO] ✅ Koneksi database berhasil ditutup")
	}()

	// Router
	header := router.InitRouter(db)

	// Log port
	log.Printf("[INFO] 🌐 Server berjalan di http://localhost:%s", port)

	// Jalankan server
	if err := http.ListenAndServe(":"+port, header); err != nil {
		log.Fatalf("[FATAL] ❌ Gagal menjalankan server: %v", err)
	}
}
