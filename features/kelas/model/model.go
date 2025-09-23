package model

import (
	"go_rest_native_sekolah/features/kelas"
	"time"
)

// Kelas merepresentasikan data kelas dalam database.
// Struktur ini digunakan untuk menyimpan informasi terkait kelas
// seperti ID kelas, nama kelas, ID guru, nama guru, waktu pembaruan, dan waktu penghapusan.
type Kelas struct {
	// ID adalah ID unik untuk setiap kelas.
	ID string `json:"id"`

	// Kelas adalah nama kelas.
	Kelas string `json:"kelas"`

	// ID_Guru adalah ID guru yang mengajar di kelas ini.
	ID_Guru string `json:"id_guru"`

	// Nama_Guru adalah nama guru yang mengajar di kelas ini.
	Nama_Guru string `json:"nama_guru"`

	// Update_At adalah waktu terakhir kali kelas ini diperbarui.
	Update_At string `json:"update_at"`

	// Delete_At adalah waktu ketika kelas ini dihapus, dapat bernilai null jika belum dihapus.
	Delete_At string `json:"delete_at"`
}

// TableName digunakan untuk mengembalikan nama tabel yang digunakan dalam database.
// Fungsi ini digunakan untuk mengatur nama tabel yang digunakan dalam database
// Nama tabel yang digunakan adalah "kelas"
// Fungsi ini digunakan untuk mengatur nama tabel yang digunakan dalam database
// Fungsi ini digunakan oleh GORM untuk mencari nama tabel yang digunakan dalam database
func (k *Kelas) TableName() string {
	// Kembalikan nama tabel yang digunakan dalam database
	return "kelas"
}

// FormatterRequest digunakan untuk mengubah objek KelasCore menjadi objek Kelas.
// Fungsi ini memformat data kelas agar sesuai dengan kebutuhan database.
func FormatterRequest(req kelas.KelasCore) Kelas {
	// Membuat objek Kelas dan mengisi dengan data dari objek KelasCore.
	return Kelas{
		ID:        req.ID,                                   // Mengisi field ID dengan ID dari objek KelasCore.
		Kelas:     req.Kelas,                                // Mengisi field Kelas dengan nama kelas dari objek KelasCore.
		ID_Guru:   req.ID_Guru,                              // Mengisi field ID_Guru dengan ID guru dari objek KelasCore.
		Nama_Guru: req.Nama_Guru,                            // Mengisi field Nama_Guru dengan nama guru dari objek KelasCore.
		Update_At: time.Now().Format("2006-01-02 15:04:05"), // Mengisi field Update_At dengan waktu saat ini.
	}
}

// FormatterResponse digunakan untuk mengubah objek Kelas menjadi objek KelasCore
// Fungsi ini digunakan untuk memformat data kelas agar sesuai dengan kebutuhan response API
// Fungsi ini menerima parameter objek Kelas dan mengembalikan objek KelasCore
func FormatterResponse(res Kelas) kelas.KelasCore {
	// Membuat objek KelasCore yang akan diisi dengan data-data kelas
	return kelas.KelasCore{
		ID:        res.ID,        // Mengisi field ID dengan ID dari objek Kelas
		Kelas:     res.Kelas,     // Mengisi field Kelas dengan nama kelas dari objek Kelas
		ID_Guru:   res.ID_Guru,   // Mengisi field ID_Guru dengan ID guru dari objek Kelas
		Nama_Guru: res.Nama_Guru, // Mengisi field Nama_Guru dengan nama guru dari objek Kelas
		Update_At: res.Update_At, // Mengisi field Update_At dengan waktu terakhir kelas diupdate dari objek Kelas
		Delete_At: res.Delete_At, // Mengisi field Delete_At dengan waktu ketika kelas dihapus dari objek Kelas
	}
}
