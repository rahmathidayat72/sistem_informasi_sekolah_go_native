package controllers

import "go_rest_native_sekolah/features/siswa"

// SiswaFormatter digunakan untuk memformat data siswa agar sesuai dengan kebutuhan response API.
// Struktur ini merepresentasikan data siswa yang akan dikirimkan sebagai respons.
type SiswaFormatter struct {
	// ID adalah field yang berisi id siswa
	ID string `json:"id"`
	// Nama adalah field yang berisi nama siswa
	Nama string `json:"nama"`
	// Kelas_ID adalah field yang berisi id kelas siswa
	Kelas_ID string `json:"kelas_id"`
	// Nama_Kelas adalah field yang berisi nama kelas siswa
	Nama_Kelas string `json:"nama_kelas"`
	// Email adalah field yang berisi email siswa
	Email string `json:"email"`
	// Alamat adalah field yang berisi alamat siswa
	Alamat string `json:"alamat"`
}

// FormatterKelasList digunakan untuk mengubah slice SiswaCore menjadi slice SiswaFormatter.
// Slice SiswaFormatter ini akan di kirimkan sebagai response API.
// Fungsi ini menerima parameter slice SiswaCore dan mengembalikan slice SiswaFormatter.
func FormatterKelasList(cores []siswa.SiswaCore) []SiswaFormatter {
	// Membuat slice SiswaFormatter yang akan diisi dengan data-data siswa
	formatted := make([]SiswaFormatter, 0)
	// Melakukan perulangan untuk setiap data siswa di dalam slice cores
	for _, core := range cores {
		// Membuat objek SiswaFormatter dan mengisi dengan data-data siswa
		formatted = append(formatted, SiswaFormatter{
			// ID adalah field yang berisi id siswa
			ID: core.ID,
			// Nama adalah field yang berisi nama siswa
			Nama: core.Nama,
			// Kelas_ID adalah field yang berisi id kelas siswa
			Kelas_ID: core.Kelas_ID,
			// Nama_Kelas adalah field yang berisi nama kelas siswa
			Nama_Kelas: core.Nama_Kelas,
			// Email adalah field yang berisi email siswa
			Email: core.Email,
			// Alamat adalah field yang berisi alamat siswa
			Alamat: core.Alamat,
		})
	}
	// Mengembalikan slice SiswaFormatter yang telah di format
	return formatted

}

// FormatSiswaRequestToCore digunakan untuk mengubah objek SiswaFormatter menjadi objek SiswaCore.
// Fungsi ini digunakan untuk memformat data siswa yang diinputkan oleh pengguna agar sesuai dengan kebutuhan aplikasi internal.
// Fungsi ini menerima parameter objek SiswaFormatter dan mengembalikan objek SiswaCore.
func FormatSiswaRequestToCore(req SiswaFormatter) siswa.SiswaCore {
	// Mengembalikan objek SiswaCore yang berisi data-data siswa dari objek SiswaFormatter
	return siswa.SiswaCore{
		ID:         req.ID,         // Mengisi field ID dengan ID dari objek SiswaFormatter
		Nama:       req.Nama,       // Mengisi field Nama dengan nama dari objek SiswaFormatter
		Kelas_ID:   req.Kelas_ID,   // Mengisi field Kelas_ID dengan ID kelas dari objek SiswaFormatter
		Nama_Kelas: req.Nama_Kelas, // Mengisi field Nama_Kelas dengan nama kelas dari objek SiswaFormatter
		Email:      req.Email,      // Mengisi field Email dengan email dari objek SiswaFormatter
		Alamat:     req.Alamat,     // Mengisi field Alamat dengan alamat dari objek SiswaFormatter
	}
}
