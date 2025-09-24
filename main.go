package main

import (
	"go_rest_native_sekolah/config"
	"go_rest_native_sekolah/router"
	"log"
	"net/http"
	"os"
)

func main() {
	// Log informasi bahwa aplikasi telah dimulai
	log.Println("[INFO] 🚀 Aplikasi dimulai...")

	// Inisialisasi koneksi database PostgreSQL
	db, err := config.InitPostgreSQLPool()
	if err != nil {
		// Jika gagal, log fatal dan hentikan aplikasi
		log.Fatalf("[FATAL] ❌ Gagal terhubung ke database: %v", err)
	}
	defer func() {
		// Tutup koneksi database saat aplikasi dihentikan
		db.Close()
		log.Println("[INFO] ✅ Koneksi database berhasil ditutup")
	}()

	// Inisialisasi router dengan koneksi database yang sudah diinisialisasi
	header := router.InitRouter(db)

	// Tentukan port untuk menjalankan server HTTP
	port := os.Getenv("PORT")
	
	// Log informasi bahwa server HTTP berjalan
	log.Printf("[INFO] 🌐 Server berjalan di http://localhost:%s", port)

	// Jalankan server HTTP dan tangani error jika terjadi
	if err := http.ListenAndServe(":"+port, header); err != nil {
		// Jika terjadi error saat menjalankan server, log fatal dan hentikan aplikasi
		log.Fatalf("[FATAL] ❌ Gagal menjalankan server: %v", err)
	}
}
