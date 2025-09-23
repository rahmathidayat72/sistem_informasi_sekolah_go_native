package model

import matapelajaran "go_rest_native_sekolah/features/mata_pelajaran"

// MataPelajaran adalah struktur yang digunakan untuk merepresentasikan data mata pelajaran.
// Struktur ini digunakan sebagai model untuk mengakses database dan memanipulasi data.
type MataPelajaran struct {
	// ID adalah field yang digunakan untuk menyimpan ID mata pelajaran.
	ID string `json:"id"`

	// Nama_Pelajaran adalah field yang digunakan untuk menyimpan nama mata pelajaran.
	Nama_Pelajaran string `json:"mata_pelajaran"`

	// ID_Guru adalah field yang digunakan untuk menyimpan ID guru yang mengajar mata pelajaran.
	ID_Guru string `json:"id_guru"`

	// Guru adalah field yang digunakan untuk menyimpan nama guru yang mengajar mata pelajaran.
	Guru string `json:"guru"`

	// Kelas_ID adalah field yang digunakan untuk menyimpan ID kelas yang mengajar mata pelajaran.
	Kelas_ID string `json:"kelas_id"`

	// Nama_Kelas adalah field yang digunakan untuk menyimpan nama kelas yang mengajar mata pelajaran.
	Nama_Kelas string `json:"nama_kelas"`

	// Deskripsi adalah field yang digunakan untuk menyimpan deskripsi mata pelajaran.
	Deskripsi string `json:"deskripsi"`

	// Update_At adalah field yang digunakan untuk menyimpan waktu terakhir data mata pelajaran diupdate.
	Update_At string `json:"update_at"`

	// Delete_At adalah field yang digunakan untuk menyimpan waktu terakhir data mata pelajaran dihapus.
	Delete_At string `json:"delete_at"`
}

// TableName adalah metode yang digunakan untuk mengembalikan nama tabel
// yang sesuai dengan struktur MataPelajaran dalam database.
// Fungsi ini mengembalikan string yang merepresentasikan nama tabel.
func (m *MataPelajaran) TableName() string {
	// Mengembalikan string "mata_pelajaran" sebagai nama tabel.
	return "mata_pelajaran"
}

// FormatterRequest adalah fungsi yang digunakan untuk mengubah objek MataPelajaranCore menjadi objek MataPelajaran.
// Fungsi ini digunakan untuk memformat data mata pelajaran yang diinputkan oleh pengguna agar sesuai dengan kebutuhan database.
// Fungsi ini menerima parameter objek MataPelajaranCore dan mengembalikan objek MataPelajaran.
func FormatterRequest(req matapelajaran.MataPelajaranCore) MataPelajaran {
	// Mengisi field ID dengan ID dari objek MataPelajaranCore
	ID := req.ID
	// Mengisi field Nama_Pelajaran dengan nama pelajaran dari objek MataPelajaranCore
	Nama_Pelajaran := req.Nama_Pelajaran
	// Mengisi field ID_Guru dengan ID guru yang mengajar mata pelajaran dari objek MataPelajaranCore
	ID_Guru := req.ID_Guru
	// Mengisi field Guru dengan nama guru yang mengajar mata pelajaran dari objek MataPelajaranCore
	Guru := req.Guru
	// Mengisi field Kelas_ID dengan ID kelas yang mengajar mata pelajaran dari objek MataPelajaranCore
	Kelas_ID := req.Kelas_ID
	// Mengisi field Nama_Kelas dengan nama kelas yang mengajar mata pelajaran dari objek MataPelajaranCore
	Nama_Kelas := req.Nama_Kelas
	// Mengisi field Deskripsi dengan deskripsi mata pelajaran dari objek MataPelajaranCore
	Deskripsi := req.Deskripsi
	// Mengisi field Update_At dengan waktu terakhir data mata pelajaran diupdate dari objek MataPelajaranCore
	Update_At := req.Update_At

	// Mengembalikan objek MataPelajaran yang telah di format
	return MataPelajaran{
		ID:             ID,
		Nama_Pelajaran: Nama_Pelajaran,
		ID_Guru:        ID_Guru,
		Guru:           Guru,
		Kelas_ID:       Kelas_ID,
		Nama_Kelas:     Nama_Kelas,
		Deskripsi:      Deskripsi,
		Update_At:      Update_At,
	}
}

// FormatterResponse adalah fungsi yang digunakan untuk mengubah objek MataPelajaran menjadi objek MataPelajaranCore.
// Fungsi ini digunakan untuk memformat data mata pelajaran yang diambil dari database agar sesuai dengan kebutuhan response API.
// Fungsi ini menerima parameter objek MataPelajaran dan mengembalikan objek MataPelajaranCore.
func FormatterResponse(res MataPelajaran) matapelajaran.MataPelajaranCore {
	// Mengisi field ID dengan ID dari objek MataPelajaran
	ID := res.ID
	// Mengisi field Nama_Pelajaran dengan nama pelajaran dari objek MataPelajaran
	Nama_Pelajaran := res.Nama_Pelajaran
	// Mengisi field ID_Guru dengan ID guru yang mengajar mata pelajaran dari objek MataPelajaran
	ID_Guru := res.ID_Guru
	// Mengisi field Guru dengan nama guru yang mengajar mata pelajaran dari objek MataPelajaran
	Guru := res.Guru
	// Mengisi field Kelas_ID dengan ID kelas yang mengajar mata pelajaran dari objek MataPelajaran
	Kelas_ID := res.Kelas_ID
	// Mengisi field Nama_Kelas dengan nama kelas yang mengajar mata pelajaran dari objek MataPelajaran
	Nama_Kelas := res.Nama_Kelas
	// Mengisi field Deskripsi dengan deskripsi mata pelajaran dari objek MataPelajaran
	Deskripsi := res.Deskripsi

	// Mengembalikan objek MataPelajaranCore yang telah di format
	return matapelajaran.MataPelajaranCore{
		ID:             ID,
		Nama_Pelajaran: Nama_Pelajaran,
		ID_Guru:        ID_Guru,
		Guru:           Guru,
		Kelas_ID:       Kelas_ID,
		Nama_Kelas:     Nama_Kelas,
		Deskripsi:      Deskripsi,
	}
}
