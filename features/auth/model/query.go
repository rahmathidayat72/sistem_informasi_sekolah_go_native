package model

import (
	"context"
	"errors"
	"go_rest_native_sekolah/features/auth"
	"go_rest_native_sekolah/helper"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AuthQuery merepresentasikan query yang berhubungan dengan autentikasi.
// Struct ini menggunakan database PostgreSQL untuk menghandle query ke database.
type AuthQuery struct {
	DB *pgxpool.Pool
}

// NewAuthData membuat objek AuthQuery dengan parameter db.
// Fungsi ini akan mengembalikan nilai AuthQuery yang siap digunakan.
// Jika parameter db nil maka akan terjadi panic.
func NewAuthData(db *pgxpool.Pool) auth.DataAuthInterface {
	if db == nil {
		// Jika db nil maka akan terjadi panic
		panic("NewAuthData: db is nil")
	}

	return &AuthQuery{DB: db}
}

// Login implements auth.DataAuthInterface.
// Fungsi ini digunakan untuk melakukan login dengan input email dan password.
// Fungsi ini akan mengembalikan nilai UserCore yang berisi data user jika login berhasil,
// atau error jika login gagal.
func (a *AuthQuery) Login(email string, password string) (dataLogin auth.UserCore, err error) {
	var userLogin User

	// Ambil user berdasarkan email saja
	// Query ini digunakan untuk mengambil data user dari database berdasarkan email.
	// Jika user tidak ditemukan maka akan terjadi error.
	query := "SELECT id, username, email, password, role FROM users WHERE email = $1"
	err = a.DB.QueryRow(context.Background(), query, email).Scan(
		&userLogin.ID,
		&userLogin.Username,
		&userLogin.Email,
		&userLogin.Password,
		&userLogin.Role,
	)

	if err != nil {
		// Jika error maka cek apakah error tersebut adalah error karena user tidak ditemukan
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("User with email %s not found", email)
			return auth.UserCore{}, errors.New("user not found")
		}

		// Jika bukan error karena user tidak ditemukan maka log error-nya
		log.Printf("Error while querying user with email %s: %v", email, err)
		return auth.UserCore{}, err
	}

	log.Printf("User found with email %s, checking password", email)

	// Bandingkan password input dengan hash password dari DB
	// Fungsi helper.CheckPassword digunakan untuk membandingkan password input dengan hash password dari DB
	// Jika password tidak sama maka akan terjadi error
	if !helper.CheckPassword(password, userLogin.Password) {
		log.Printf("Login failed for user with email %s: wrong password", email)
		return auth.UserCore{}, errors.New("login failed, wrong password")
	}

	log.Printf("Login successful for user with email %s", email)

	// Buatkan objek UserCore berdasarkan data user yang diambil dari database
	dataLogin = auth.UserCore(FormatterResponse(userLogin))
	return dataLogin, nil
}
