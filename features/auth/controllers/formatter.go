package controllers

import "time"

type (
	// ResponseAuth digunakan untuk merepresentasikan respons setelah login berhasil.
	// Struktur ini berisi ID pengguna, nama pengguna, email, token autentikasi, dan waktu kedaluwarsa token.
	ResponseAuth struct {
		ID         string    `json:"id"`         // ID pengguna
		Username   string    `json:"username"`   // Nama pengguna
		Email      string    `json:"email"`      // Email pengguna
		Token      string    `json:"token"`      // Token autentikasi yang diberikan setelah login berhasil
		Expiration time.Time `json:"expiration"` // Waktu kedaluwarsa token
	}

	// LoginRequest digunakan untuk merepresentasikan permintaan login.
	// Struktur ini berisi email dan password yang digunakan untuk login.
	LoginRequest struct {
		Email    string `json:"email"`    // Email pengguna untuk proses login
		Password string `json:"password"` // Password pengguna untuk proses login
	}
)
