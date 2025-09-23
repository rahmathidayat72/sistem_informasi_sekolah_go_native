package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/auth"
	"go_rest_native_sekolah/helper"
	"log"
	"net/http"
	"strings"
)

// AuthController merepresentasikan controller untuk autentikasi.
// authService digunakan untuk menghandle logika bisnis yang berhubungan dengan autentikasi.
type AuthController struct {
	authService auth.ServiceAuthInterface
}

// NewAutController membuat objek AuthController dengan parameter authService.
// Fungsi ini digunakan untuk menghandle logika bisnis terkait autentikasi.
// Jika parameter authService nil, maka akan terjadi panic.
func NewAutController(authService auth.ServiceAuthInterface) *AuthController {
	// Mengembalikan objek AuthController baru yang siap digunakan
	return &AuthController{
		authService: authService, // Menyimpan service autentikasi ke dalam field authService
	}
}

// Auth melakukan login user berdasarkan input yang diterima.
// Fungsi ini akan mengembalikan token dan data user jika login berhasil.
// Jika login gagal maka akan terjadi error.
func (lc *AuthController) Auth(w http.ResponseWriter, r *http.Request) error {
	if lc.authService == nil {
		return errors.New("auth controller: Nil service")
	}

	var inputLogin LoginRequest // Input yang diterima dari request body
	if err := json.NewDecoder(r.Body).Decode(&inputLogin); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Data tidak valid", http.StatusBadRequest)
		return fmt.Errorf("auth controller: error decoding request: %v", err)
	}

	log.Printf("Attempting to log in user with email: %s", inputLogin.Email)
	login, err := lc.authService.Login(inputLogin.Email, inputLogin.Password) // Melakukan login
	if err != nil {
		log.Printf("Login error for %s: %v", inputLogin.Email, err)
		if strings.Contains(err.Error(), "invalid Password") {
			http.Error(w, "Email atau password salah", http.StatusUnauthorized)
		} else {
			http.Error(w, "Terjadi kesalahan saat login", http.StatusInternalServerError)
		}
		return err
	}

	data := map[string]interface{}{"id": login.ID} // Data untuk membuat token
	token, expTime, err := helper.SignToken(data)  // Membuat token berdasarkan data
	if err != nil {
		log.Printf("Token generation failed: %v", err)
		http.Error(w, "Gagal membuat token", http.StatusInternalServerError)
		return err
	}

	// Membuat response yang berisi token, id, username, email, dan expiration time
	response := ResponseAuth{
		ID:         login.ID,
		Username:   login.Username, // atau Nama, konsisten
		Email:      login.Email,
		Token:      token,
		Expiration: expTime,
	}

	w.Header().Set("Content-Type", "application/json") // Menentukan tipe konten response agar dapat di parse oleh client-side
	// Izinkan request dari semua origin untuk keperluan development saja
	// Dalam production, pastikan untuk mengatur Access-Control-Allow-Origin
	// agar hanya mengizinkan request dari domain yang terpercaya
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Membuat response yang berisi kode status, message, dan data
	apiResponse := helper.APIResponse(http.StatusOK, "success login", response) // Membuat response
	if err := json.NewEncoder(w).Encode(apiResponse); err != nil {              // Mengencode response
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError) // Mengirimkan response error
		return err
	}

	log.Printf("Login successful for user with email: %s", inputLogin.Email)
	return nil // Mengembalikan nil karena login berhasil
}
