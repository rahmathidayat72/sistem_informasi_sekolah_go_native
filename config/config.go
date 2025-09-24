package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv berfungsi untuk load file .env sesuai mode
func LoadEnv() {
	// Cek environment APP_ENV
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development" // default kalau tidak di-set
	}

	envFile := "dev.env" // Ganti bagian ini ketika ada file env lain atau diubah ke production.env

	// Coba load file sesuai APP_ENV
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("[WARN] Tidak menemukan %s, mencoba .env default", envFile)

		// fallback ke .env default
		if err := godotenv.Load("exp.env"); err != nil {
			log.Printf("[WARN] Tidak menemukan .env default, menggunakan environment bawaan OS")
		}
	} else {
		log.Printf("[INFO] ✅ Berhasil load %s", envFile)
	}
}
