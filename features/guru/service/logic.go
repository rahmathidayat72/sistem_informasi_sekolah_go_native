package service

import (
	"context"
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/guru"
	"regexp"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// guruService  merepresentasikan service untuk tabel guru
type guruService struct {
	guruData guru.DataGuruInterface // guruData  berisi kumpulan function-pointers yang dibutuhkan untuk mengakses data guru
	db       *pgxpool.Pool          // db berisi kumpulan function-pointers yang dibutuhkan untuk mengakses database
}

// SelectById implements guru.ServiceGuruInterface.

// InsertGuru implements guru.ServiceGuruInterface.

// NewServiceGuru digunakan untuk membuat objek guruService dengan parameter guruData.
// guruService digunakan untuk menghandle logika bisnis yang berhubungan dengan tabel guru.
// Jika parameter guruData nil maka akan terjadi panic.
func NewServiceGuru(repo guru.DataGuruInterface, db *pgxpool.Pool) guru.ServiceGuruInterface {
	if repo == nil {
		panic("guru service: Nil repository")
	}
	return &guruService{guruData: repo,
		db: db}

}

// GetAllGuru digunakan untuk mengambil semua data guru dari database.
// Fungsi ini akan mengembalikan slice guru.GuruCore yang berisi data guru
// dan error jika terjadi kesalahan.
func (s *guruService) GetAllGuru() ([]guru.GuruCore, error) {
	// Periksa apakah guruData adalah nil
	if s.guruData == nil {
		// Kembalikan error jika guruData nil
		return nil, errors.New("guru service: Nil repository")
	}

	// Panggil fungsi SelectAllGuru dari guruData untuk mengambil semua data guru
	gurus, err := s.guruData.SelectAllGuru()
	if err != nil {
		// Kembalikan error jika terjadi kesalahan saat mengambil data
		return nil, errors.New("guru service: gagal mengambil data")
	}

	// Kembalikan slice dari guru.GuruCore dan error nil
	return gurus, nil
}

// InsertGuru digunakan untuk memasukkan data guru ke dalam database.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (s *guruService) InsertGuru(insert *guru.GuruCore) error {
	// Periksa apakah data repository guruData kosong (nil).
	if s.guruData == nil {
		return errors.New("guru service: Repository kosong")
	}

	// Validasi bahwa nama, email, dan alamat tidak boleh kosong.
	if insert.Nama == "" || insert.Email == "" || insert.Alamat == "" {
		return errors.New("validasi error: nama, email, dan alamat harus diisi")
	}

	// Validasi format email menggunakan regex.
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(insert.Email) {
		return errors.New("validasi error: email tidak valid")
	}

	// Deklarasi variabel untuk menyimpan ID user.
	var idUser string

	// Cek apakah email sudah ada di tabel users untuk mendapatkan id_user.
	err := s.db.QueryRow(context.Background(),
		"SELECT id FROM users WHERE email = $1", insert.Email).Scan(&idUser)

	if err != nil {
		// Jika email tidak ditemukan, biarkan id_user kosong (NULL di database).
		if err == pgx.ErrNoRows {
			idUser = ""
		} else {
			// Jika ada kesalahan lain saat pengecekan, kembalikan error.
			return fmt.Errorf("gagal cek users: %w", err)
		}
	}

	// Panggil fungsi InsertGuru pada data repository untuk memasukkan data guru.
	return s.guruData.InsertGuru(insert)
}

// UpdateGuru memperbarui data guru berdasarkan ID yang diberikan.
// Fungsi ini mengimplementasikan guru.ServiceGuruInterface.
func (s *guruService) UpdateGuru(insert *guru.GuruCore, id string) error {
	// Periksa apakah service atau data repository nil
	if s == nil || s.guruData == nil {
		return errors.New("guru service: Nil repository")
	}

	// Periksa apakah input data nil
	if insert == nil {
		return errors.New("guru service: input is nil")
	}

	// Validasi ID harus diisi
	if id == "" {
		return errors.New("validation error: id harus diisi")
	}

	// Ambil data lama dari database berdasarkan ID
	existingData, err := s.guruData.SelectById(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Jika data tidak ditemukan, kembalikan error
			return errors.New("guru service: Data tidak ditemukan")
		}
		return fmt.Errorf("Id salah atau gagal mengambil data lama: %w", err)
	}

	// Gabungkan data baru dengan data lama
	// Jika field baru kosong, gunakan field dari data lama
	if insert.Nama == "" {
		insert.Nama = existingData.Nama
	}
	if insert.Email == "" {
		insert.Email = existingData.Email
	}
	if insert.Alamat == "" {
		insert.Alamat = existingData.Alamat
	}

	// Lakukan update data ke database
	if err := s.guruData.Update(insert, id); err != nil {
		return err
	}

	// Kembalikan nil jika berhasil
	return nil
}

// SelectById digunakan untuk mengambil data guru berdasarkan ID yang diberikan.
// Fungsi ini akan mengembalikan objek guru.GuruCore yang sesuai dengan ID tersebut
// dan error jika terjadi kesalahan dalam pengambilan data.
func (s *guruService) SelectById(id string) (*guru.GuruCore, error) {
	// Panggil method SelectById dari data layer untuk mengambil data guru berdasarkan ID.
	guru, err := s.guruData.SelectById(id)
	// Jika terjadi error saat mengambil data, kembalikan nil dan error yang terbungkus dengan pesan.
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data guru: %w", err)
	}
	// Kembalikan objek guru yang berhasil diambil dan error nil.
	return guru, nil
}

// DeleteById mengimplementasikan interface guru.ServiceGuruInterface.
// Fungsi ini digunakan untuk menghapus data guru berdasarkan ID yang diberikan.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (s *guruService) DeleteById(id string) error {
	// Periksa apakah service atau data repository nil.
	if s == nil || s.guruData == nil {
		return errors.New("guru service: Nil repository")
	}

	// Periksa apakah ID harus diisi.
	if id == "" {
		return errors.New("validation error: id harus diisi")
	}

	// Panggil fungsi DeleteById pada data repository untuk menghapus data guru.
	// Jika terjadi error saat menghapus data guru, kembalikan error yang terbungkus dengan pesan.
	if err := s.guruData.DeleteById(id); err != nil {
		return fmt.Errorf("gagal menghapus data guru: %w", err)
	}

	// Jika tidak ada error maka kembalikan nil.
	return nil
}
