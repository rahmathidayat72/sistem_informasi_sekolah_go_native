package model

import (
	"go_rest_native_sekolah/features/users"
	"time"
)

// User adalah struktur data yang merepresentasikan data user
// dalam database.
type User struct {
	// ID adalah field yang berisi id user yang unik
	ID string `json:"id"`
	// Username adalah field yang berisi nama user yang unik
	Username string `json:"username"`
	// Email adalah field yang berisi email user yang unik
	Email string `json:"email"`
	// Password adalah field yang berisi password user yang dienkripsi
	Password string `json:"password"`
	// Role adalah field yang berisi peran user dalam sistem
	// peran user dapat berupa "admin" atau "user"
	Role string `json:"role"`
	// Update_at adalah field yang berisi waktu terakhir update data user
	Update_At string `json:"update_at"`
	// Delete_at adalah field yang berisi waktu penghapusan data user
	// jika field ini kosong maka data user tidak pernah dihapus
	Delete_At string `json:"delete_at"`
}

// TableName mengembalikan nama tabel yang terkait dengan struktur data User.
// Fungsi ini digunakan oleh ORM atau query builder untuk menentukan nama tabel di database.
func (u *User) TableName() string {
	// Mengembalikan string "users" sebagai nama tabel yang sesuai dengan struktur User di database.
	return "users"
}

// FormatterRequest adalah fungsi yang digunakan untuk mengubah objek UserCore menjadi objek User.
// Fungsi ini digunakan untuk memformat data user agar sesuai dengan kebutuhan database.
// Fungsi ini menerima parameter objek UserCore dan mengembalikan objek User.
func FormatterRequest(req users.UserCore) User {
	// Membuat objek User dan mengisi dengan data dari objek UserCore.
	return User{
		ID:        req.ID,                                   // Mengisi field ID dengan ID dari objek UserCore.
		Username:  req.Username,                             // Mengisi field Username dengan Username dari objek UserCore.
		Email:     req.Email,                                // Mengisi field Email dengan Email dari objek UserCore.
		Password:  req.Password,                             // Mengisi field Password dengan Password dari objek UserCore.
		Role:      req.Role,                                 // Mengisi field Role dengan Role dari objek UserCore.
		Update_At: time.Now().Format("2006-01-02 15:04:05"), // Mengisi field Update_At dengan waktu saat ini.
	}

}

// FormatterResponse adalah fungsi yang digunakan untuk mengubah objek User menjadi objek UserCore.
// Fungsi ini digunakan untuk memformat data user agar sesuai dengan kebutuhan aplikasi internal.
// Fungsi ini menerima parameter objek User dan mengembalikan objek UserCore.
func FormatterResponse(res User) users.UserCore {
	// Mengembalikan objek UserCore yang berisi data user dari objek User.
	return users.UserCore{
		ID:       res.ID,       // Mengisi field ID dengan ID dari objek User.
		Username: res.Username, // Mengisi field Username dengan Username dari objek User.
		Email:    res.Email,    // Mengisi field Email dengan Email dari objek User.
		Password: res.Password, // Mengisi field Password dengan Password dari objek User.
		Role:     res.Role,     // Mengisi field Role dengan Role dari objek User.
	}
}
