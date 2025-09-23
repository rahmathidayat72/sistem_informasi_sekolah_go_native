package controllers

import "go_rest_native_sekolah/features/users"

// UserFormatter adalah struktur data yang merepresentasikan data user yang akan dikirimkan sebagai respon API.
// Struktur ini berisi ID user, nama pengguna, email, password, dan role user.
type UserFormatter struct {
	ID       string `json:"id"`       // ID adalah identifikasi unik untuk setiap user.
	Username string `json:"username"` // Username adalah nama pengguna yang digunakan untuk login.
	Email    string `json:"email"`    // Email adalah alamat email user yang digunakan untuk login.
	Password string `json:"password"` // Password adalah password yang digunakan user untuk login.
	Role     string `json:"role"`     // Role adalah peran user yang menentukan akses terhadap fitur-fitur di aplikasi.

}

// FormatUserList adalah fungsi yang digunakan untuk mengubah slice UserCore menjadi slice UserFormatter.
// Slice UserFormatter ini akan di kirimkan sebagai response API.
// Fungsi ini menerima parameter slice UserCore dan mengembalikan slice UserFormatter.
func FormatUserList(cores []users.UserCore) []UserFormatter {
	// Membuat slice UserFormatter yang akan diisi dengan data-data user
	formatted := make([]UserFormatter, 0)
	// Melakukan perulangan untuk setiap data user di dalam slice cores
	for _, core := range cores {
		// Membuat objek UserFormatter dan mengisi dengan data-data user
		formatted = append(formatted, UserFormatter{
			ID:       core.ID,       // ID adalah identifikasi unik untuk setiap user.
			Username: core.Username, // Username adalah nama pengguna yang digunakan untuk login.
			Email:    core.Email,    // Email adalah alamat email user yang digunakan untuk login.
			Password: core.Password, // Password adalah password yang digunakan user untuk login.
			Role:     core.Role,     // Role adalah peran user yang menentukan akses terhadap fitur-fitur di aplikasi.
		})
	}
	// Mengembalikan slice UserFormatter yang telah di format
	return formatted

}

// FormatUserRequestToCore adalah fungsi yang digunakan untuk mengubah objek UserFormatter menjadi objek UserCore.
// Fungsi ini menerima parameter objek UserFormatter dan mengembalikan objek UserCore.
// Fungsi ini digunakan untuk memformat data user yang diinputkan oleh pengguna agar sesuai dengan kebutuhan database.
func FormatUserRequestToCore(req UserFormatter) users.UserCore {
	// Membuat objek UserCore dan mengisi dengan data-data user
	return users.UserCore{
		ID:       req.ID,       // ID adalah identifikasi unik untuk setiap user.
		Username: req.Username, // Username adalah nama pengguna yang digunakan untuk login.
		Email:    req.Email,    // Email adalah alamat email user yang digunakan untuk login.
		Password: req.Password, // Password adalah password yang digunakan user untuk login.
		Role:     req.Role,     // Role adalah peran user yang menentukan akses terhadap fitur-fitur di aplikasi.
	}
}
