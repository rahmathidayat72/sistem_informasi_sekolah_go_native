package controllers

import "go_rest_native_sekolah/features/guru"

// GuruFormatter digunakan untuk memformat data guru agar sesuai dengan kebutuhan response API.
// Struktur ini merepresentasikan data guru yang akan dikirimkan sebagai respons.
type GuruFormatter struct {
	ID      string `json:"id"`      // ID adalah ID unik untuk setiap guru
	ID_User string `json:"id_user"` // ID_User adalah ID dari pengguna yang terkait dengan guru
	Nama    string `json:"nama"`    // Nama adalah nama lengkap dari guru
	Email   string `json:"email"`   // Email adalah alamat email dari guru
	Alamat  string `json:"alamat"`  // Alamat adalah alamat tempat tinggal dari guru
}

// FormatGuruList digunakan untuk mengubah slice GuruCore menjadi slice GuruFormatter.
// Slice GuruFormatter ini akan di kirimkan sebagai response API.
// Fungsi ini menerima parameter slice GuruCore dan mengembalikan slice GuruFormatter.
func FormatGuruList(cores []guru.GuruCore) []GuruFormatter {
	formatted := make([]GuruFormatter, 0)
	// Looping slice GuruCore dan mengubah setiap elemen menjadi GuruFormatter.
	// Kemudian, append GuruFormatter ke dalam slice formatted.
	for _, core := range cores {
		formatted = append(formatted, GuruFormatter{
			ID:      core.ID,      // ID adalah ID unik untuk setiap guru
			ID_User: core.ID_User, // ID_User adalah ID dari pengguna yang terkait dengan guru
			Nama:    core.Nama,    // Nama adalah nama lengkap dari guru
			Email:   core.Email,   // Email adalah alamat email dari guru
			Alamat:  core.Alamat,  // Alamat adalah alamat tempat tinggal dari guru
		})
	}
	// Mengembalikan slice GuruFormatter yang telah di format.
	return formatted
}

// FormatGuruRequestToCore digunakan untuk mengubah objek GuruFormatter menjadi objek GuruCore.
// Fungsi ini digunakan untuk memformat data guru yang diinputkan oleh pengguna agar sesuai dengan kebutuhan database.
// Fungsi ini mengembalikan objek GuruCore yang berisi data guru yang telah di format.
func FormatGuruRequestToCore(req GuruFormatter) guru.GuruCore {
	// Membuat objek GuruCore yang akan diisi dengan data guru yang diinputkan.
	core := guru.GuruCore{}
	// ID adalah ID unik untuk setiap guru.
	core.ID = req.ID
	// ID_User adalah ID dari pengguna yang terkait dengan guru.
	core.ID_User = req.ID_User
	// Nama adalah nama lengkap dari guru.
	core.Nama = req.Nama
	// Email adalah alamat email dari guru.
	core.Email = req.Email
	// Alamat adalah alamat tempat tinggal dari guru.
	core.Alamat = req.Alamat
	// Mengembalikan objek GuruCore yang telah di format.
	return core
}
