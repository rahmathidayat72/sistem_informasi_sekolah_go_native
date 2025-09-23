package controllers

import (
	"go_rest_native_sekolah/features/kelas"
	"time"
)

// KelasFormatter digunakan untuk memformat data kelas agar sesuai dengan kebutuhan response API.
// Struktur ini merepresentasikan data kelas yang akan dikirimkan sebagai respons.
type KelasFormatter struct {
	// ID adalah ID unik untuk setiap kelas
	ID string `json:"id"`
	// Kelas adalah nama kelas
	Kelas string `json:"kelas"`
	// ID_Guru adalah ID guru yang mengajar di kelas ini
	ID_Guru string `json:"id_guru"`
	// Nama_Guru adalah nama guru yang mengajar di kelas ini
	Nama_Guru string `json:"nama_guru"`
}

// FormatKelasList digunakan untuk mengubah slice KelasCore menjadi slice KelasFormatter.
// Slice KelasFormatter ini akan di kirimkan sebagai response API.
// Fungsi ini menerima parameter slice KelasCore dan mengembalikan slice KelasFormatter.
func FormatKelasList(cores []kelas.KelasCore) []KelasFormatter {
	// Membuat slice KelasFormatter yang akan diisi dengan data-data kelas
	formatted := make([]KelasFormatter, 0)
	// Melakukan perulangan untuk setiap data kelas di dalam slice cores
	for _, core := range cores {
		// Membuat objek KelasFormatter dan mengisi dengan data-data kelas
		formatted = append(formatted, KelasFormatter{
			// ID adalah ID unik untuk setiap kelas
			ID: core.ID,
			// Kelas adalah nama kelas
			Kelas: core.Kelas,
			// ID_Guru adalah ID guru yang mengajar di kelas ini
			ID_Guru: core.ID_Guru,
			// Nama_Guru adalah nama guru yang mengajar di kelas ini
			Nama_Guru: core.Nama_Guru,
		})
	}
	// Mengembalikan slice KelasFormatter yang telah di format
	return formatted

}

// FormatKelasRequestToCore digunakan untuk mengubah objek KelasFormatter menjadi objek KelasCore.
// Fungsi ini digunakan untuk memformat data kelas yang diinputkan oleh pengguna agar sesuai dengan kebutuhan database.
// Fungsi ini menerima parameter objek KelasFormatter dan mengembalikan objek KelasCore.
func FormatKelasRequestToCore(req KelasFormatter) kelas.KelasCore {
	// Membuat objek KelasCore yang akan diisi dengan data-data kelas
	core := kelas.KelasCore{}
	// ID adalah ID unik untuk setiap kelas
	core.ID = req.ID
	// Kelas adalah nama kelas
	core.Kelas = req.Kelas
	// ID_Guru adalah ID guru yang mengajar di kelas ini
	core.ID_Guru = req.ID_Guru
	// Nama_Guru adalah nama guru yang mengajar di kelas ini
	core.Nama_Guru = req.Nama_Guru
	// Update_At adalah waktu terakhir kelas diupdate
	core.Update_At = time.Now().Format("2006-01-02 15:04:05")
	// Mengembalikan objek KelasCore yang telah di format
	return core
}
