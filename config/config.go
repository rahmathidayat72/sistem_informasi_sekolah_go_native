package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv berfungsi untuk memuat variabel lingkungan dari file sesuai dengan APP_ENV.
// Fungsi ini akan:
// - Mengecek variabel lingkungan APP_ENV untuk menentukan file .env yang akan dimuat.
// - Jika APP_ENV tidak di-set, maka default ke "development" dan memuat file ".dev.env".
// - Untuk nilai APP_ENV lainnya, akan memuat file dengan format ".<APP_ENV>.env" (misal: ".production.env").
// - Jika file yang spesifik tidak ditemukan, akan fallback ke file ".env".
// - Jika kedua file tidak ditemukan, maka menggunakan variabel lingkungan dari OS.
// Fungsi ini juga akan mencatat informasi dan peringatan terkait proses pemuatan file .env.
func LoadEnv() {
	// Cek environment APP_ENV
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development" // default kalau tidak di-set
	}

	// File .env yang akan dimuat sesuai dengan APP_ENV
	envFile := fmt.Sprintf(".%s.env", env)
	if env == "development" {
		envFile = ".dev.env" // Bisa di set sesuai kebutuhan dan file env yang digunakan
	}

	// Coba load file sesuai APP_ENV
	if err := godotenv.Load(envFile); err != nil {
		// Jika tidak menemukan file .env yang spesifik, maka mencoba .env default
		log.Printf("[WARN] ❌ Tidak menemukan %s, mencoba .env default", envFile)

		// fallback ke .env (global/default)
		if err := godotenv.Load(".env"); err != nil {
			// Jika tidak menemukan .env default, maka menggunakan environment bawaan OS
			log.Printf("[WARN] ❌ Tidak menemukan .env default, menggunakan environment bawaan OS")
		} else {
			log.Printf("[INFO] ✅ Berhasil load .env default")
		}
	} else {
		log.Printf("[INFO] ✅ Berhasil load %s", envFile)
	}
}
