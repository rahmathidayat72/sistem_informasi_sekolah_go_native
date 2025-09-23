package model

import (
	"go_rest_native_sekolah/features/users"
)

// User merepresentasikan data user di database
// Data user terdiri atas:
// 1. ID (string) sebagai id data user
// 2. Username (string) sebagai nama pengguna
// 3. Email (string) sebagai alamat email user
// 4. Password (string) sebagai password user
// 5. Role (string) sebagai peran user
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// TableName digunakan untuk mengembalikan nama tabel yang digunakan
// oleh struct User. Nama tabel yang digunakan adalah "users".
func (u *User) TableName() string {
	return "users"

}

// FormatterRequest digunakan untuk mengubah objek UserCore menjadi objek User
// Fungsi ini digunakan untuk memformat data user agar sesuai dengan kebutuhan database
func FormatterRequest(req users.UserCore) User {
	return User{
		Email:    req.Email,
		Password: req.Password,
	}

}

// FormatterResponse digunakan untuk mengubah objek User menjadi objek UserCore
// Fungsi ini digunakan untuk memformat data user agar sesuai dengan kebutuhan response API
func FormatterResponse(res User) users.UserCore {
	return users.UserCore{
		ID:       res.ID,
		Username: res.Username,
		Email:    res.Email,
		Password: res.Password,
		Role:     res.Role,
	}
}
