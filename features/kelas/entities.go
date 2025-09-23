package kelas

// KelasCore adalah struct yang merepresentasikan data kelas di database
// Struktur ini digunakan untuk menyimpan informasi terkait kelas
// seperti ID kelas, nama kelas, ID guru, nama guru, waktu terakhir diperbarui, dan waktu dihapus
type KelasCore struct {
	ID        string `json:"id"`        // ID kelas
	Kelas     string `json:"kelas"`     // Nama kelas
	ID_Guru   string `json:"id_guru"`   // ID guru
	Nama_Guru string `json:"nama_guru"` // Nama guru
	Update_At string `json:"update_at"` // Waktu terakhir diperbarui
	Delete_At string `json:"delete_at"` // Waktu dihapus
}

// DataKelasInterface adalah interface yang berhubungan dengan data kelas
// Interface ini memiliki method SelectAll, SelectById, Insert, Update, dan DeleteById
// Method-method ini digunakan untuk menghandle data kelas di database
type DataKelasInterface interface {
	// SelectAll digunakan untuk mengambil semua data kelas di database
	// Fungsi ini mengembalikan slice KelasCore yang berisi data kelas
	// Jika terjadi error maka fungsi ini akan mengembalikan error
	SelectAll() ([]KelasCore, error)
	// SelectById digunakan untuk mengambil data kelas berdasarkan ID yang diberikan
	// Fungsi ini akan mengembalikan objek KelasCore yang sesuai dengan ID tersebut
	// dan error jika terjadi kesalahan dalam pengambilan data
	SelectById(id string) (*KelasCore, error)
	// Insert digunakan untuk menginsert data kelas ke dalam database
	// Fungsi ini mengembalikan error jika terjadi kesalahan
	Insert(insert *KelasCore) error
	// Update digunakan untuk mengupdate data kelas berdasarkan ID yang diberikan
	// Fungsi ini akan mengembalikan error jika terjadi kesalahan dalam proses update
	Update(insert *KelasCore, id string) error
	// DeleteById digunakan untuk menghapus data kelas berdasarkan ID yang diberikan
	// Fungsi ini akan mengembalikan error jika terjadi kesalahan dalam proses hapus
	DeleteById(id string) error
}

// ServiceKelasInterface adalah interface yang berhubungan dengan service kelas
// Interface ini memiliki method SelectAll, SelectById, Insert, Update, dan DeleteById
// Method-method ini digunakan untuk menghandle data kelas di database
type ServiceKelasInterface interface {
	// SelectAll digunakan untuk mengambil semua data kelas di database
	// Fungsi ini mengembalikan slice KelasCore yang berisi data kelas
	// Jika terjadi error maka fungsi ini akan mengembalikan error
	SelectAll() ([]KelasCore, error)
	// SelectById digunakan untuk mengambil data kelas berdasarkan ID yang diberikan
	// Fungsi ini akan mengembalikan objek KelasCore yang sesuai dengan ID tersebut
	// dan error jika terjadi kesalahan dalam pengambilan data
	SelectById(id string) (*KelasCore, error)
	// Insert digunakan untuk menginsert data kelas ke dalam database
	// Fungsi ini mengembalikan error jika terjadi kesalahan
	Insert(insert *KelasCore) error
	// Update digunakan untuk mengupdate data kelas berdasarkan ID yang diberikan
	// Fungsi ini akan mengembalikan error jika terjadi kesalahan dalam proses update
	Update(insert *KelasCore, id string) error
	// DeleteById digunakan untuk menghapus data kelas berdasarkan ID yang diberikan
	// Fungsi ini akan mengembalikan error jika terjadi kesalahan dalam proses hapus
	DeleteById(id string) error
}
