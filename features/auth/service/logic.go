package service

import (
	"go_rest_native_sekolah/features/auth"
	"log"
)

// authService merepresentasikan service untuk autentikasi.
// authService digunakan untuk menghandle logika bisnis yang berhubungan dengan autentikasi.
// authService berisi field authData yang digunakan untuk mengakses data autentikasi.
type authService struct {
	authData auth.DataAuthInterface // authData digunakan untuk mengakses data autentikasi.
}

// NewServiceAuth membuat objek authService yang siap digunakan.
// authService digunakan untuk menghandle logika bisnis yang berhubungan dengan autentikasi.
// Jika parameter authData nil maka akan terjadi panic.
func NewServiceAuth(authData auth.DataAuthInterface) auth.ServiceAuthInterface {
	if authData == nil {
		panic("NewServiceAuth: authData is nil")
	}
	// Membuat objek authService yang siap digunakan
	return &authService{
		authData: authData, // Menyimpan data auth ke dalam field authData
	}
}

// Login implements auth.ServiceAuthInterface.
// Fungsi ini digunakan untuk melakukan login dengan input email dan password.
// Fungsi ini akan mengembalikan nilai UserCore yang berisi data user jika login berhasil,
// atau error jika login gagal.
// Fungsi ini juga akan mengembalikan error jika terjadi kesalahan saat login.
func (a *authService) Login(email string, password string) (dataLogin auth.UserCore, err error) {
	log.Printf("Memulai proses login untuk email: %s", email)
	// Menggunakan data auth untuk menghandle login
	// Jika terjadi error maka akan mengembalikan error
	dataLogin, err = a.authData.Login(email, password)
	if err != nil {
		log.Printf("Terjadi kesalahan saat login user dengan email %s: %v", email, err)
		return auth.UserCore{}, err // Mengembalikan error jika terjadi kesalahan
	}
	log.Printf("Login berhasil untuk email: %s", email)
	// Mengembalikan data user yang berhasil login
	return dataLogin, nil
}
