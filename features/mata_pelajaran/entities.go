package matapelajaran

// MataPelajaranCore adalah struktur data yang berisi field2 yang akan diisi
// oleh data mata pelajaran.
type MataPelajaranCore struct {
	// ID adalah field yang berisi ID unik untuk setiap mata pelajaran.
	ID string `json:"id"`
	// Nama_Pelajaran adalah field yang berisi nama mata pelajaran.
	Nama_Pelajaran string `json:"mata_pelajaran"`
	// ID_Guru adalah field yang berisi ID guru yang mengajar mata pelajaran.
	ID_Guru string `json:"id_guru"`
	// Guru adalah field yang berisi nama guru yang mengajar mata pelajaran.
	Guru string `json:"guru"`
	// Kelas_ID adalah field yang berisi ID kelas yang mengajar mata pelajaran.
	Kelas_ID string `json:"kelas_id"`
	// Nama_Kelas adalah field yang berisi nama kelas yang mengajar mata pelajaran.
	Nama_Kelas string `json:"nama_kelas"`
	// Deskripsi adalah field yang berisi deskripsi singkat tentang mata pelajaran.
	Deskripsi string `json:"deskripsi"`
	// Update_At adalah field yang berisi waktu update terakhir data mata pelajaran.
	Update_At string `json:"update_at"`
	// Delete_At adalah field yang berisi waktu delete data mata pelajaran.
	Delete_At string `json:"delete_at"`
}

// DataMataPelajaranInterface adalah interface yang berisi method2 yang digunakan
// untuk mengambil data mata pelajaran dari database dan melakukan operasi CRUD.
type DataMataPelajaranInterface interface {
	// SelectAllMapel adalah method yang digunakan untuk mengambil semua data mata pelajaran
	// dari database.
	SelectAllMapel() ([]MataPelajaranCore, error)
	// SelectMapelById adalah method yang digunakan untuk mengambil data mata pelajaran
	// berdasarkan ID dari database.
	SelectMapelById(id string) (*MataPelajaranCore, error)
	// InsertMapel adalah method yang digunakan untuk menginsert data mata pelajaran
	// ke dalam database.
	InsertMapel(insert *MataPelajaranCore) error
	// UpdateMapel adalah method yang digunakan untuk mengupdate data mata pelajaran
	// berdasarkan ID di database.
	UpdateMapel(insert *MataPelajaranCore, id string) error
	// DeleteMapel adalah method yang digunakan untuk menghapus data mata pelajaran
	// berdasarkan ID di database.
	DeleteMapel(id string) error
}

// ServiceMapelInterface adalah interface yang berisi method2 yang digunakan
// untuk menghandle request dari client dan mengoperasikan data mata pelajaran.
type ServiceMapelInterface interface {
	// SelectAllMapel adalah method yang digunakan untuk mengambil semua data mata pelajaran
	// dari database.
	SelectAllMapel() ([]MataPelajaranCore, error)
	// SelectMapelById adalah method yang digunakan untuk mengambil data mata pelajaran
	// berdasarkan ID dari database.
	SelectMapelById(id string) (*MataPelajaranCore, error)
	// InsertMapel adalah method yang digunakan untuk menginsert data mata pelajaran
	// ke dalam database.
	InsertMapel(insert *MataPelajaranCore) error
	// UpdateMapel adalah method yang digunakan untuk mengupdate data mata pelajaran
	// berdasarkan ID di database.
	UpdateMapel(insert *MataPelajaranCore, id string) error
	// DeleteMapel adalah method yang digunakan untuk menghapus data mata pelajaran
	// berdasarkan ID di database.
	DeleteMapel(id string) error
}
