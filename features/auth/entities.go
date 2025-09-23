package auth

import "time"

type (
	// UserCore merepresentasikan data user di database.
	// Data user terdiri atas:
	// 1. ID (string) sebagai id data user
	// 2. Username (string) sebagai nama pengguna
	// 3. Email (string) sebagai alamat email user
	// 4. Password (string) sebagai password user
	// 5. Role (string) sebagai peran user
	// 6. Update_At (time.Time) sebagai waktu update data user
	// 7. Delete_At (*time.Time) sebagai waktu delete data user
	UserCore struct {
		ID        string     `json:"id"`        // ID data user
		Username  string     `json:"username"`  // Nama pengguna
		Email     string     `json:"email"`     // Email user
		Password  string     `json:"password"`  // Password user
		Role      string     `json:"role"`      // Peran user
		Update_At time.Time  `json:"update_at"` // Waktu update data user
		Delete_At *time.Time `json:"delete_at"` // Waktu delete data user
	}

	// DataAuthInterface merepresentasikan interface untuk data auth.
	// Interface ini digunakan untuk menghandle data auth yang berhubungan dengan user.
	DataAuthInterface interface {
		// Login melakukan login user berdasarkan input email dan password.
		// Fungsi ini akan mengembalikan nilai UserCore yang berisi data user jika login berhasil,
		// atau error jika login gagal.
		Login(email, password string) (dataLogin UserCore, err error)
	}

	// ServiceAuthInterface merepresentasikan interface untuk service auth.
	// Interface ini digunakan untuk menghandle logika bisnis yang berhubungan dengan autentikasi.
	ServiceAuthInterface interface {
		// Login melakukan login user berdasarkan input email dan password.
		// Fungsi ini akan mengembalikan nilai UserCore yang berisi data user jika login berhasil,
		// atau error jika login gagal.
		Login(email, password string) (dataLogin UserCore, err error)
	}
)
