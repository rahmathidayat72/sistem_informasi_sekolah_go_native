package model

import (
	"go_rest_native_sekolah/features/siswa"
	"time"
)

// Siswa adalah struktur data yang merepresentasikan informasi siswa.
// Struktur ini digunakan untuk menyimpan informasi siswa yang ada dalam database.
type Siswa struct {
	ID         string `json:"id"`         // ID adalah identifikasi unik untuk setiap siswa.
	Kelas_ID   string `json:"kelas_id"`   // Kelas_ID adalah ID dari kelas tempat siswa berada.
	Nama       string `json:"nama"`       // Nama adalah nama lengkap siswa.
	Nama_Kelas string `json:"nama_kelas"` // Nama_Kelas adalah nama kelas tempat siswa berada.
	Email      string `json:"email"`      // Email adalah alamat email siswa.
	Alamat     string `json:"alamat"`     // Alamat adalah alamat tempat tinggal siswa.
	Update_At  string `json:"update_at"`  // Update_At adalah waktu terakhir data siswa diperbarui.
	Delete_At  string `json:"delete_at"`  // Delete_At adalah waktu ketika data siswa dihapus, jika ada.
}

// TableName mengembalikan nama tabel yang terkait dengan struktur data Siswa.
// Fungsi ini digunakan oleh ORM atau query builder untuk menentukan nama tabel di database.
func (u *Siswa) TableName() string {
	// Mengembalikan string "siswa" sebagai nama tabel yang sesuai dengan struktur Siswa di database.
	return "siswa"
}

// FormatterRequest adalah fungsi yang digunakan untuk mengubah objek SiswaCore menjadi objek Siswa.
// Fungsi ini digunakan untuk memformat data siswa yang diinputkan oleh pengguna agar sesuai dengan kebutuhan database.
// Fungsi ini menerima parameter objek SiswaCore dan mengembalikan objek Siswa.
func FormatterRequest(req siswa.SiswaCore) Siswa {
	// Membuat objek Siswa yang akan diisi dengan data siswa yang diinputkan.
	siswa := Siswa{}
	// ID adalah identifikasi unik untuk setiap siswa.
	siswa.ID = req.ID
	// Kelas_ID adalah ID dari kelas tempat siswa berada.
	siswa.Kelas_ID = req.Kelas_ID
	// Nama adalah nama lengkap siswa.
	siswa.Nama = req.Nama
	// Nama_Kelas adalah nama kelas tempat siswa berada.
	siswa.Nama_Kelas = req.Nama_Kelas
	// Email adalah alamat email siswa.
	siswa.Email = req.Email
	// Alamat adalah alamat tempat tinggal siswa.
	siswa.Alamat = req.Alamat
	// Update_At adalah waktu terakhir data siswa diperbarui.
	siswa.Update_At = time.Now().Format("2006-01-02 15:04:05")
	// Mengembalikan objek Siswa yang telah di format.
	return siswa
}

// FormatterResponse adalah fungsi yang digunakan untuk mengubah objek Siswa menjadi objek SiswaCore.
// Fungsi ini digunakan untuk memformat data siswa agar sesuai dengan kebutuhan aplikasi internal.
// Fungsi ini menerima parameter objek Siswa dan mengembalikan objek SiswaCore.
func FormatterResponse(res Siswa) siswa.SiswaCore {
	// Mengembalikan objek SiswaCore yang berisi data-data siswa dari objek Siswa
	return siswa.SiswaCore{
		ID:         res.ID,         // Mengisi field ID dengan ID dari objek Siswa
		Nama:       res.Nama,       // Mengisi field Nama dengan nama dari objek Siswa
		Kelas_ID:   res.Kelas_ID,   // Mengisi field Kelas_ID dengan ID kelas dari objek Siswa
		Nama_Kelas: res.Nama_Kelas, // Mengisi field Nama_Kelas dengan nama kelas dari objek Siswa
		Email:      res.Email,      // Mengisi field Email dengan email dari objek Siswa
		Alamat:     res.Alamat,     // Mengisi field Alamat dengan alamat dari objek Siswa
	}
}
