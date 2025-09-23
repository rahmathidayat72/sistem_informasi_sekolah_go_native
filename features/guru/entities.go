package guru

import "time"

type ( // GuruCore struct untuk merepresentasikan tabel guru
	GuruCore struct { // Guru struct untuk merepresentasikan tabel guru
		ID        string     `json:"id"`
		ID_User   string     `json:"id_user"`
		Nama      string     `json:"nama"`
		Email     string     `json:"email"`
		Alamat    string     `json:"alamat"`
		Update_At time.Time  `json:"update_at"`
		Delete_At *time.Time `json:"delete_at"`
	}

	DataGuruInterface interface { // Interface untuk mengakses data guru
		// SelectAllGuru digunakan untuk mengambil semua data guru dari database.
		// Fungsi ini mengembalikan slice dari GuruCore yang berisi data guru.
		// Jika terjadi kesalahan selama pengambilan data, fungsi ini akan mengembalikan error.
		SelectAllGuru() ([]GuruCore, error)
		InsertGuru(insert *GuruCore) error
		Update(insert *GuruCore, id string) error
		SelectById(id string) (*GuruCore, error)
		DeleteById(id string) error
	}

	ServiceGuruInterface interface { // Interface untuk mengakses logika bisnis guru
		// GetAllGuru digunakan untuk mengambil semua data guru dari database.
		// Fungsi ini mengembalikan slice dari GuruCore yang berisi data guru.
		// Jika terjadi kesalahan selama pengambilan data, fungsi ini akan mengembalikan error.
		GetAllGuru() ([]GuruCore, error)
		InsertGuru(insert *GuruCore) error
		UpdateGuru(insert *GuruCore, id string) error
		SelectById(id string) (*GuruCore, error)
		DeleteById(id string) error
	}
)
