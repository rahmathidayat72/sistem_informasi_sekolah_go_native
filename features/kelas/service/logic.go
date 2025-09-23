package service

import (
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/kelas"

	"github.com/jackc/pgx/v5"
)

// kelasService merepresentasikan service yang berhubungan dengan data kelas.
// Struct ini memiliki satu field yaitu kelasData yang berisi interface DataKelasInterface.
// kelasData digunakan untuk mengakses data kelas dari repository.
type kelasService struct {
	kelasData kelas.DataKelasInterface // Menyimpan referensi ke interface DataKelasInterface
}

// NewServiceKelas digunakan untuk membuat objek service kelas yang berhubungan dengan data kelas.
// Fungsi ini memiliki parameter repo yang berisi interface DataKelasInterface.
// Parameter repo digunakan untuk mengakses data kelas dari repository.
// Jika parameter repo nil maka akan terjadi panic.
func NewServiceKelas(repo kelas.DataKelasInterface) kelas.ServiceKelasInterface {
	// Cek apakah parameter repo nil
	if repo == nil {
		// Jika parameter repo nil maka akan terjadi Error
		errors.New("Nil repository")
	}

	// Membuat objek service kelas dengan parameter repo
	return &kelasService{kelasData: repo}
}

// Insert implements kelas.ServiceKelasInterface.
// Fungsi ini digunakan untuk menginsert data kelas ke dalam database.
// Fungsi ini memiliki parameter insert yang berisi objek kelas.KelasCore.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan saat menginsert data.
func (k *kelasService) Insert(insert *kelas.KelasCore) error {
	// Memeriksa apakah koneksi ke repository ada atau tidak
	if k.kelasData == nil {
		// Jika koneksi repository nil maka akan terjadi panic
		panic("Nil repository")
	}

	// Memeriksa apakah parameter insert nil
	if insert.Kelas == "" {
		// Jika parameter insert nil maka akan terjadi error
		return errors.New("Validation error: insert kelas is nil")
	}

	// Menginsert data kelas ke dalam database menggunakan repository
	return k.kelasData.Insert(insert)
}

// SelectAll digunakan untuk mengambil semua data kelas dari repository.
// Fungsi ini mengembalikan slice KelasCore dan error jika terjadi kesalahan.
func (k *kelasService) SelectAll() ([]kelas.KelasCore, error) {
	// Memeriksa apakah koneksi ke repository ada atau tidak
	if k.kelasData == nil {
		// Jika repository nil, kembalikan error
		return nil, errors.New("kelas service: Nil repository")
	}

	// Mengambil semua data kelas dari repository
	kelass, err := k.kelasData.SelectAll()
	if err != nil {
		// Jika terjadi error saat pengambilan data, kembalikan error
		return nil, errors.New("kelas service: gagal mengambil data")
	}

	// Mengembalikan data kelas yang berhasil diambil
	return kelass, nil
}

// SelectById implements kelas.ServiceKelasInterface.
// Fungsi ini digunakan untuk mengambil data kelas berdasarkan ID yang diberikan.
// Fungsi ini mengembalikan objek kelas.KelasCore yang sesuai dengan ID tersebut
// dan error jika terjadi kesalahan dalam pengambilan data.
func (k *kelasService) SelectById(id string) (*kelas.KelasCore, error) {
	// Memeriksa apakah koneksi ke repository ada atau tidak
	// Jika repository nil, kembalikan error
	if k.kelasData == nil {
		return nil, errors.New("kelas service: Nil repository")
	}

	// Mengambil data kelas berdasarkan ID yang diberikan
	// Fungsi ini akan mengembalikan error jika terjadi kesalahan
	// dalam pengambilan data.
	kelas, err := k.kelasData.SelectById(id)
	if err != nil {
		// Jika terjadi error saat pengambilan data, kembalikan error
		// dengan pesan yang sesuai.
		return nil, fmt.Errorf("kelas service: gagal mengambil data kelas %w", err)
	}

	// Mengembalikan data kelas yang berhasil diambil
	return kelas, nil
}

// Update digunakan untuk memperbarui data kelas berdasarkan ID yang diberikan.
// Fungsi ini mengembalikan error jika terjadi kesalahan dalam proses update.
func (k *kelasService) Update(insert *kelas.KelasCore, id string) error {
	// Memeriksa apakah service atau data repository nil
	if k == nil || k.kelasData == nil {
		return errors.New("Nil repository")
	}

	// Memeriksa apakah ID kosong
	if id == "" {
		return errors.New("Validation error: id is nil")
	}

	// Mengambil data kelas yang ada berdasarkan ID
	exisData, err := k.kelasData.SelectById(id)
	if err != nil {
		// Jika data tidak ditemukan, kembalikan error
		if err == pgx.ErrNoRows {
			return errors.New("Data tidak ditemukan")
		}
		// Kembalikan error jika terjadi kesalahan lain saat mengambil data
		return fmt.Errorf("Id salah atau gagal mengambil data lama: %w", err)
	}

	// Gabungkan data baru dengan data lama
	// Jika field baru kosong, gunakan field dari data lama
	if insert.Kelas == "" {
		insert.Kelas = exisData.Kelas
	}
	if insert.Nama_Guru == "" {
		insert.Nama_Guru = exisData.Nama_Guru
	}
	if insert.ID_Guru == "" {
		insert.ID_Guru = exisData.ID_Guru
	}

	// Memperbarui data kelas ke dalam database
	if err := k.kelasData.Update(insert, id); err != nil {
		// Kembalikan error jika terjadi kesalahan saat memperbarui data
		return err
	}

	// Kembalikan nil jika update berhasil tanpa error
	return nil
}

// DeleteById implements kelas.ServiceKelasInterface.
// Fungsi ini digunakan untuk menghapus data kelas berdasarkan ID yang diberikan.
// Fungsi ini memiliki parameter id yang berisi string ID kelas yang ingin dihapus.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan saat menghapus data.
func (k *kelasService) DeleteById(id string) error {
	// Memeriksa apakah service atau data repository nil
	if k == nil || k.kelasData == nil {
		// Jika repository nil, kembalikan error
		return errors.New("Nil repository")
	}

	// Memeriksa apakah ID kosong
	if id == "" {
		// Jika ID kosong, kembalikan error
		return errors.New("Validation error: id harus diisi")
	}

	// Menghapus data kelas berdasarkan ID yang diberikan
	// Fungsi ini akan mengembalikan error jika terjadi kesalahan
	// dalam penghapusan data.
	if err := k.kelasData.DeleteById(id); err != nil {
		// Jika terjadi error saat menghapus data, kembalikan error
		// dengan pesan yang sesuai.
		return fmt.Errorf("gagal menghapus data kelas: %w", err)
	}

	// Kembalikan nil jika hapus berhasil tanpa error
	return nil
}
