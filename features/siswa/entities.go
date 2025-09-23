package siswa

import "time"

type (
	// SiswaCore adalah struktur data yang merepresentasikan informasi inti dari seorang siswa.
	// Struktur ini berisi ID siswa, nama, ID kelas, nama kelas, email, alamat, dan informasi waktu pembaruan serta penghapusan.
	SiswaCore struct {
		ID         string     `json:"id"`         // ID adalah identifikasi unik untuk setiap siswa.
		Nama       string     `json:"nama"`       // Nama adalah nama lengkap siswa.
		Kelas_ID   string     `json:"kelas_id"`   // Kelas_ID adalah ID dari kelas tempat siswa berada.
		Nama_Kelas string     `json:"nama_kelas"` // Nama_Kelas adalah nama kelas tempat siswa berada.
		Email      string     `json:"email"`      // Email adalah alamat email siswa.
		Alamat     string     `json:"alamat"`     // Alamat adalah alamat tempat tinggal siswa.
		Update_At  time.Time  `json:"update_at"`  // Update_At adalah waktu terakhir data siswa diperbarui.
		Delete_At  *time.Time `json:"delete_at"`  // Delete_At adalah waktu di mana data siswa dihapus, jika ada.
	}

	// DataSiswaInterface adalah antarmuka yang mendefinisikan metode untuk operasi data siswa.
	// Antarmuka ini mencakup metode untuk mengambil semua data siswa, memasukkan data siswa,
	// memperbarui data siswa, mengambil data siswa berdasarkan ID, dan menghapus data siswa berdasarkan ID.
	DataSiswaInterface interface {
		SelectAllSiswa() ([]SiswaCore, error)      // Mengambil semua data siswa dari database.
		InsertSiswa(insert *SiswaCore) error       // Memasukkan data siswa baru ke dalam database.
		Update(insert *SiswaCore, id string) error // Memperbarui data siswa berdasarkan ID.
		SelectById(id string) (*SiswaCore, error)  // Mengambil data siswa berdasarkan ID.
		DeleteById(id string) error                // Menghapus data siswa berdasarkan ID.
	}

	// ServiceSiswaInterface adalah antarmuka yang mendefinisikan layanan untuk operasi siswa.
	// Antarmuka ini serupa dengan DataSiswaInterface, namun digunakan di lapisan layanan untuk
	// mengabstraksi operasi-operasi yang dilakukan pada data siswa.
	ServiceSiswaInterface interface {
		SelectAllSiswa() ([]SiswaCore, error)      // Mengambil semua data siswa dari database.
		InsertSiswa(insert *SiswaCore) error       // Memasukkan data siswa baru ke dalam database.
		Update(insert *SiswaCore, id string) error // Memperbarui data siswa berdasarkan ID.
		SelectById(id string) (*SiswaCore, error)  // Mengambil data siswa berdasarkan ID.
		DeleteById(id string) error                // Menghapus data siswa berdasarkan ID.
	}
)
