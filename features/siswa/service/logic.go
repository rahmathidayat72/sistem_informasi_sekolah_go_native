package service

import (
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/siswa"
	"regexp"

	"github.com/jackc/pgx/v5"
)

// siswaService adalah struct yang digunakan untuk mengimplementasikan interface ServiceSiswaInterface.
// Struct ini memiliki satu field yaitu siswaData yang merupakan interface DataSiswaInterface.
// siswaData digunakan untuk mengakses data siswa dari database.
type siswaService struct {
	siswaData siswa.DataSiswaInterface // Interface untuk mengakses data siswa dari database
}

// NewServiceSiswa digunakan untuk membuat objek siswaService yang akan digunakan
// untuk menghandle logika bisnis yang berhubungan dengan data siswa.
// Fungsi ini menerima parameter repo yang berupa interface DataSiswaInterface.
// Parameter repo digunakan untuk mengakses data siswa dari database.
// Jika parameter repo nil maka akan terjadi panic.
func NewServiceSiswa(repo siswa.DataSiswaInterface) siswa.ServiceSiswaInterface {
	// Cek apakah parameter repo nil atau tidak.
	// Jika nil maka akan terjadi panic.
	if repo == nil {
		panic("Nil repository")
	}
	// Membuat objek siswaService yang berisi parameter repo.
	// Objek siswaService ini akan digunakan untuk menghandle logika bisnis
	// yang berhubungan dengan data siswa.
	return &siswaService{siswaData: repo}
}

// InsertSiswa implements siswa.ServiceSiswaInterface.
// InsertSiswa digunakan untuk memasukkan data siswa ke dalam database.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (s *siswaService) InsertSiswa(insert *siswa.SiswaCore) error {
	// Memeriksa apakah repository siswaData tidak nil.
	if s.siswaData == nil {
		return errors.New("Nil repository")
	}
	// Memeriksa apakah nama siswa tidak kosong.
	if insert.Nama == "" {
		return errors.New("Validation error: nama siswa tidak boleh kosong")
	}
	// Memeriksa apakah alamat siswa tidak kosong.
	if insert.Alamat == "" {
		return errors.New("Validation error: alamat siswa tidak boleh kosong")
	}
	// Memeriksa apakah email siswa tidak kosong.
	if insert.Email == "" {
		return errors.New("Validation error: email siswa tidak boleh kosong")
	}
	// Menyiapkan regex untuk memvalidasi format email.
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	// Memeriksa apakah email sesuai dengan format yang benar.
	if !emailRegex.MatchString(insert.Email) {
		return errors.New("validation error: email tidak valid")
	}

	// Memanggil fungsi InsertSiswa pada siswaData untuk menyimpan data siswa.
	return s.siswaData.InsertSiswa(insert)
}

// SelectAllSiswa implements siswa.ServiceSiswaInterface.
// SelectAllSiswa implements siswa.ServiceSiswaInterface.
// Fungsi ini digunakan untuk mengambil semua data siswa dari database.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (s *siswaService) SelectAllSiswa() ([]siswa.SiswaCore, error) {
	// Memeriksa apakah repository siswaData tidak nil.
	if s.siswaData == nil {
		return nil, errors.New("SiswaService: Nil repository")
	}
	// Memanggil fungsi SelectAllSiswa pada siswaData untuk mengambil data siswa.
	kelass, err := s.siswaData.SelectAllSiswa()
	// Jika terjadi error maka kembalikan error.
	if err != nil {
		return nil, errors.New("SiswaService: gagal mengambil data")
	}
	// Kembalikan data siswa yang telah diambil.
	return kelass, nil
}

// SelectById implements siswa.ServiceSiswaInterface.
// Fungsi ini digunakan untuk mengambil data siswa berdasarkan ID.
// Fungsi ini akan mengembalikan data siswa yang sesuai dengan ID yang dikirimkan
// dan error jika terjadi kesalahan.
func (s *siswaService) SelectById(id string) (*siswa.SiswaCore, error) {
	// Memanggil fungsi SelectById pada siswaData untuk mengambil data siswa berdasarkan ID.
	siswa, err := s.siswaData.SelectById(id)
	if err != nil {
		// Jika terjadi error saat mengambil data siswa, mengembalikan error dengan pesan "SiswaService: gagal mengambil data".
		return nil, errors.New("SiswaService: gagal mengambil data")
	}
	// Mengembalikan data siswa yang berhasil diambil dan nil jika tidak ada error.
	return siswa, nil
}

// Update implements siswa.ServiceSiswaInterface.
// Fungsi ini digunakan untuk memperbarui data siswa berdasarkan ID.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (s *siswaService) Update(insert *siswa.SiswaCore, id string) error {
	// Memeriksa apakah repository siswaData tidak nil.
	if s == nil || s.siswaData == nil {
		return errors.New("Nil repository")
	}
	// Memeriksa apakah ID tidak kosong.
	if id == "" {
		return errors.New("Validation error: id is nil")
	}
	// Mengambil data siswa yang akan diupdate berdasarkan ID.
	existingData, err := s.siswaData.SelectById(id)
	if err != nil {
		// Jika terjadi error saat mengambil data siswa, kembalikan error.
		if err == pgx.ErrNoRows {
			// Jika data tidak ditemukan, kembalikan error.
			return errors.New("guru service: Data tidak ditemukan")
		}
		return fmt.Errorf("Id salah atau gagal mengambil data lama: %w", err)
	}
	// Menggabungkan data lama dengan data baru.
	// Jika field baru kosong, gunakan field dari data lama.
	if insert.Nama == "" {
		insert.Nama = existingData.Nama
	}
	if insert.Email == "" {
		insert.Email = existingData.Email
	}
	if insert.Alamat == "" {
		insert.Alamat = existingData.Alamat
	}
	if insert.Kelas_ID == "" {
		insert.Kelas_ID = existingData.Kelas_ID
	}
	// Memanggil fungsi Update pada siswaData untuk memperbarui data siswa.
	if err := s.siswaData.Update(insert, id); err != nil {
		// Jika terjadi error saat memperbarui data siswa, kembalikan error.
		return fmt.Errorf("failed to update data: %w", err)
	}
	// Kembalikan nil jika berhasil memperbarui data siswa.
	return nil
}

// DeleteById implements siswa.ServiceSiswaInterface.
// DeleteById implements siswa.ServiceSiswaInterface.
// Fungsi ini digunakan untuk menghapus data siswa berdasarkan ID yang dikirimkan.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (s *siswaService) DeleteById(id string) error {
	if s == nil || s.siswaData == nil {
		// Jika koneksi database tidak ada, maka kembalikan error.
		return errors.New("Nil repository")
	}
	if id == "" {
		// Jika parameter id kosong, maka kembalikan error.
		return errors.New("validation error: id harus diisi")
	}
	if err := s.siswaData.DeleteById(id); err != nil {
		// Jika terjadi error saat menghapus data siswa, kembalikan error.
		return fmt.Errorf("gagal menghapus data siswa: %w", err)
	}

	return nil // Kembalikan nil jika berhasil menghapus data siswa
}
