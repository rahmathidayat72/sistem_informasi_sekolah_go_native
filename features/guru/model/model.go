package gurumodels

import (
	"go_rest_native_sekolah/features/guru"
	"time"
)

// Guru merepresentasikan data guru di database
// Data Guru terdiri atas:
// 1. ID (int) sebagai id data guru
// 2. Nama (string) sebagai nama guru
// 3. Email (string) sebagai email guru
// 4. Alamat (string) sebagai alamat guru
// 5. Update_At (time.Time) sebagai waktu update data guru
// 6. Delete_At (*time.Time) sebagai waktu delete data guru
type Guru struct {
	ID        string     `json:"id"`
	ID_User   string     `json:"id_user"`
	Nama      string     `json:"nama"`
	Email     string     `json:"email"`
	Alamat    string     `json:"alamat"`
	Update_At time.Time  `json:"update_at"`
	Delete_At *time.Time `json:"delete_at"`
}

// TableName digunakan untuk mengembalikan nama tabel yang digunakan dalam database
// Fungsi ini digunakan untuk mengatur nama tabel yang digunakan dalam database
// Nama tabel yang digunakan adalah "guru"
func (u *Guru) TableName() string {
	return "guru" // Nama tabel yang digunakan
}

// FormatterRequest digunakan untuk mengubah objek GuruCore menjadi objek Guru.
// Fungsi ini memformat data guru agar sesuai dengan kebutuhan database.
func FormatterRequest(req guru.GuruCore) Guru {
	// Membuat objek Guru dan mengisi dengan data dari objek GuruCore.
	return Guru{
		ID:        req.ID,      // Mengisi field ID dengan ID dari objek GuruCore.
		ID_User:   req.ID_User, // Mengisi field ID_User dengan ID_User dari objek GuruCore.
		Nama:      req.Nama,    // Mengisi field Nama dengan Nama dari objek GuruCore.
		Email:     req.Email,   // Mengisi field Email dengan Email dari objek GuruCore.
		Alamat:    req.Alamat,  // Mengisi field Alamat dengan Alamat dari objek GuruCore.
		Update_At: time.Now(),  // Mengisi field Update_At dengan waktu saat ini.
	}
}

// FormatterResponse digunakan untuk mengubah objek Guru menjadi objek GuruCore
// Fungsi ini digunakan untuk memformat data guru agar sesuai dengan kebutuhan response API
// Fungsi ini akan mengembalikan objek GuruCore yang berisi data guru.
func FormatterResponse(res Guru) guru.GuruCore {
	// Membuat objek GuruCore dan mengisi dengan data guru yang diambil dari database
	// ID adalah ID unik untuk setiap guru.
	// ID_User adalah ID dari pengguna yang terkait dengan guru.
	// Nama adalah nama lengkap dari guru.
	// Email adalah alamat email dari guru.
	// Alamat adalah alamat tempat tinggal dari guru.
	return guru.GuruCore{
		ID:      res.ID,
		ID_User: res.ID_User,
		Nama:    res.Nama,
		Email:   res.Email,
		Alamat:  res.Alamat,
	}
}

