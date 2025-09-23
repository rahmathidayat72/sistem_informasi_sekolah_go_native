package controllers

import matapelajaran "go_rest_native_sekolah/features/mata_pelajaran"

// FormatterMataPelajaran digunakan untuk memformat data mata pelajaran agar sesuai dengan kebutuhan response API.
// Struktur ini merepresentasikan data mata pelajaran yang akan dikirimkan sebagai respons.
type FormatterMataPelajaran struct {
	ID             string `json:"id"`             // ID adalah ID unik untuk setiap mata pelajaran
	Nama_Pelajaran string `json:"mata_pelajaran"` // Nama_Pelajaran adalah nama dari mata pelajaran
	ID_Guru        string `json:"id_guru"`        // ID_Guru adalah ID dari guru yang mengajar mata pelajaran ini
	Guru           string `json:"guru"`           // Guru adalah nama dari guru yang mengajar mata pelajaran ini
	Kelas_ID       string `json:"kelas_id"`       // Kelas_ID adalah ID dari kelas tempat mata pelajaran ini diajarkan
	Nama_Kelas     string `json:"nama_kelas"`     // Nama_Kelas adalah nama dari kelas tempat mata pelajaran ini diajarkan
	Deskripsi      string `json:"deskripsi"`      // Deskripsi adalah penjelasan singkat tentang mata pelajaran ini
}

// FormatterMapelList digunakan untuk mengubah slice MataPelajaranCore menjadi slice FormatterMataPelajaran.
// Slice FormatterMataPelajaran ini akan di kirimkan sebagai response API.
// Fungsi ini menerima parameter slice MataPelajaranCore dan mengembalikan slice FormatterMataPelajaran.
func FormatterMapelList(cores []matapelajaran.MataPelajaranCore) []FormatterMataPelajaran {
	formatted := make([]FormatterMataPelajaran, 0) // Membuat slice FormatterMataPelajaran yang kosong untuk diisi dengan data-data mata pelajaran
	for _, core := range cores {                   // Melakukan perulangan untuk setiap data mata pelajaran di dalam slice cores
		formatted = append(formatted, FormatterMataPelajaran{ // Membuat objek FormatterMataPelajaran dan mengisi dengan data-data mata pelajaran
			ID:             core.ID,             // ID adalah ID unik untuk setiap mata pelajaran
			Nama_Pelajaran: core.Nama_Pelajaran, // Nama_Pelajaran adalah nama dari mata pelajaran
			ID_Guru:        core.ID_Guru,        // ID_Guru adalah ID dari guru yang mengajar mata pelajaran ini
			Guru:           core.Guru,           // Guru adalah nama dari guru yang mengajar mata pelajaran ini
			Kelas_ID:       core.Kelas_ID,       // Kelas_ID adalah ID dari kelas tempat mata pelajaran ini diajarkan
			Nama_Kelas:     core.Nama_Kelas,     // Nama_Kelas adalah nama dari kelas tempat mata pelajaran ini diajarkan
			Deskripsi:      core.Deskripsi,      // Deskripsi adalah penjelasan singkat tentang mata pelajaran ini
		})
	}
	return formatted // Mengembalikan slice FormatterMataPelajaran yang telah di format
}

// FormatterMapelRequestToCore digunakan untuk mengubah objek FormatterMataPelajaran menjadi objek MataPelajaranCore.
// Fungsi ini memformat data mata pelajaran yang diinputkan oleh pengguna agar sesuai dengan kebutuhan aplikasi internal.
// Fungsi ini menerima parameter objek FormatterMataPelajaran dan mengembalikan objek MataPelajaranCore.
func FormatterMapelRequestToCore(req FormatterMataPelajaran) matapelajaran.MataPelajaranCore {
	// Mengembalikan objek MataPelajaranCore yang berisi data-data mata pelajaran dari objek FormatterMataPelajaran
	return matapelajaran.MataPelajaranCore{
		ID:             req.ID,             // Mengisi field ID dengan ID dari objek FormatterMataPelajaran
		Nama_Pelajaran: req.Nama_Pelajaran, // Mengisi field Nama_Pelajaran dengan nama pelajaran dari objek FormatterMataPelajaran
		ID_Guru:        req.ID_Guru,        // Mengisi field ID_Guru dengan ID guru dari objek FormatterMataPelajaran
		Guru:           req.Guru,           // Mengisi field Guru dengan nama guru dari objek FormatterMataPelajaran
		Kelas_ID:       req.Kelas_ID,       // Mengisi field Kelas_ID dengan ID kelas dari objek FormatterMataPelajaran
		Nama_Kelas:     req.Nama_Kelas,     // Mengisi field Nama_Kelas dengan nama kelas dari objek FormatterMataPelajaran
		Deskripsi:      req.Deskripsi,      // Mengisi field Deskripsi dengan deskripsi dari objek FormatterMataPelajaran
	}
}
