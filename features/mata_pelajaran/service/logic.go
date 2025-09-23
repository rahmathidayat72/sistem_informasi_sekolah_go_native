package service

import (
	"errors"
	"fmt"
	matapelajaran "go_rest_native_sekolah/features/mata_pelajaran"

	"github.com/jackc/pgx/v5"
)

// mataPelajaranServiceinterface adalah struct yang berisi field mataPelajaranData yang
// berjenis DataMataPelajaranInterface. Struct ini digunakan sebagai implementasi dari
// interface ServiceMapelInterface.
type mataPelajaranServiceinterface struct {
	mataPelajaranData matapelajaran.DataMataPelajaranInterface // field mataPelajaranData berisi pointer ke DataMataPelajaranInterface.
	// DataMataPelajaranInterface adalah interface yang berisi method-method untuk menghandle
	// query ke database yang berhubungan dengan tabel mata_pelajaran.
}

// NewMataPelajaranService adalah fungsi yang digunakan untuk membuat
// instance dari MataPelajaranServiceInterface. Fungsi ini menerima
// parameter berupa DataMataPelajaranInterface yang nantinya akan
// diisi ke dalam field mataPelajaranData di dalam struct
// mataPelajaranServiceinterface.
//
// Jika parameter yang diinputkan adalah nil maka akan terjadi Eror.
// Jika parameter yang diinputkan bukan nil maka akan dibuatkan
// instance dari MataPelajaranServiceInterface yang berisi pointer
// ke DataMataPelajaranInterface.
func NewMataPelajaranService(repo matapelajaran.DataMataPelajaranInterface) matapelajaran.ServiceMapelInterface {
	if repo == nil {
		// Jika parameter yang diinputkan adalah nil maka akan terjadi panic.
		errors.New("Nil repository")
	}
	// Membuatkan instance dari MataPelajaranServiceInterface yang berisi pointer
	// ke DataMataPelajaranInterface.
	return &mataPelajaranServiceinterface{mataPelajaranData: repo}
}

// InsertMapel digunakan untuk menginsert data mata pelajaran ke dalam database.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan saat proses insert.
func (m *mataPelajaranServiceinterface) InsertMapel(insert *matapelajaran.MataPelajaranCore) error {
	// Memeriksa apakah mataPelajaranData adalah nil.
	// Jika nil, kembalikan error karena repository tidak dapat diakses.
	if m.mataPelajaranData == nil {
		return errors.New("Nil repository")
	}

	// Memeriksa apakah data yang akan diinsert adalah nil.
	// Jika nil, kembalikan error karena tidak ada data yang bisa diproses.
	if insert == nil {
		return errors.New("insert data is nil")
	}

	// Memeriksa apakah nama pelajaran dalam data yang akan diinsert kosong.
	// Jika kosong, kembalikan error karena nama pelajaran wajib diisi.
	if insert.Nama_Pelajaran == "" {
		return errors.New("Validation error: insert mapel is nil")
	}

	// Memanggil fungsi InsertMapel pada mataPelajaranData untuk memasukkan data ke dalam database.
	// Jika terjadi error saat proses insert, error tersebut akan diteruskan.
	return m.mataPelajaranData.InsertMapel(insert)
}

// SelectAllMapel implements matapelajaran.ServiceMapelInterface.
// Fungsi ini digunakan untuk mengambil semua data mata pelajaran yang tersedia di database.
// Fungsi ini akan mengembalikan slice of MataPelajaranCore yang berisi data mata pelajaran.
// Jika terjadi error maka fungsi ini akan mengembalikan error.
func (m *mataPelajaranServiceinterface) SelectAllMapel() ([]matapelajaran.MataPelajaranCore, error) {
	// Memeriksa apakah mataPelajaranData adalah nil.
	// Jika nil, kembalikan error karena repository tidak dapat diakses.
	if m.mataPelajaranData == nil {
		return nil, errors.New("Nil repository")
	}

	// Memanggil fungsi SelectAllMapel pada mataPelajaranData untuk mengambil data.
	// Jika terjadi error saat mengambil data, error tersebut akan diteruskan.
	mapels, err := m.mataPelajaranData.SelectAllMapel()
	if err != nil {
		return nil, errors.New("MataPelajaranService: gagal mengambil data")
	}

	// Mengembalikan slice of MataPelajaranCore yang berisi data mata pelajaran.
	return mapels, nil
}

// SelectMapelById implements matapelajaran.ServiceMapelInterface.
// Fungsi ini digunakan untuk mengambil data mata pelajaran berdasarkan ID.
// Fungsi ini akan mengembalikan pointer ke struct MataPelajaranCore yang berisi data mata pelajaran.
// Jika terjadi error maka fungsi ini akan mengembalikan error.
func (m *mataPelajaranServiceinterface) SelectMapelById(id string) (*matapelajaran.MataPelajaranCore, error) {
	// Memeriksa apakah mataPelajaranData adalah nil.
	// Jika nil, kembalikan error karena repository tidak dapat diakses.
	if m.mataPelajaranData == nil {
		return nil, errors.New("Nil repository")
	}

	// Memanggil fungsi SelectMapelById pada mataPelajaranData untuk mengambil data berdasarkan ID.
	// Jika terjadi error saat mengambil data, error tersebut akan diteruskan.
	mapel, err := m.mataPelajaranData.SelectMapelById(id)
	if err != nil {
		return nil, errors.New("MataPelajaranService: gagal mengambil data")
	}

	// Mengembalikan pointer ke struct MataPelajaranCore yang berisi data mata pelajaran.
	return mapel, nil
}

// UpdateMapel implements matapelajaran.ServiceMapelInterface.
func (m *mataPelajaranServiceinterface) UpdateMapel(update *matapelajaran.MataPelajaranCore, id string) error {
	if m == nil || m.mataPelajaranData == nil {
		return errors.New("Nil repository")
	}

	if id == "" {
		return errors.New("Validation error: id is nil")
	}

	// Ambil data lama berdasarkan ID
	existingData, err := m.mataPelajaranData.SelectMapelById(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("mata pelajaran service: Data tidak ditemukan")
		}
		return fmt.Errorf("gagal mengambil data lama: %w", err)
	}

	// Merge data jika field baru kosong
	if update.Nama_Pelajaran == "" {
		update.Nama_Pelajaran = existingData.Nama_Pelajaran
	}
	if update.ID_Guru == "" {
		update.ID_Guru = existingData.ID_Guru
	}
	if update.Kelas_ID == "" {
		update.Kelas_ID = existingData.Kelas_ID
	}
	if update.Deskripsi == "" {
		update.Deskripsi = existingData.Deskripsi
	}

	// Lakukan update ke database
	if err := m.mataPelajaranData.UpdateMapel(update, id); err != nil {
		return fmt.Errorf("gagal update data mata pelajaran: %w", err)
	}

	return nil // Mengembalikan nil jika update berhasil
}

// DeleteMapel implements matapelajaran.ServiceMapelInterface.
// Fungsi ini digunakan untuk menghapus data mata pelajaran berdasarkan ID.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan saat proses hapus.
func (m *mataPelajaranServiceinterface) DeleteMapel(id string) error {
	// Memeriksa apakah mataPelajaranData adalah nil.
	// Jika nil, kembalikan error karena repository tidak dapat diakses.
	if m == nil || m.mataPelajaranData == nil {
		return errors.New("Nil repository")
	}
	// Memeriksa apakah ID yang akan dihapus kosong.
	// Jika kosong, kembalikan error karena ID wajib diisi.
	if id == "" {
		return errors.New("Validation error: id is nil")
	}
	// Memanggil fungsi DeleteMapel pada mataPelajaranData untuk menghapus data berdasarkan ID.
	// Jika terjadi error saat proses hapus, error tersebut akan diteruskan.
	if err := m.mataPelajaranData.DeleteMapel(id); err != nil {
		// Jika terjadi error maka kembalikan error dengan pesan "gagal menghapus data mata pelajaran".
		return fmt.Errorf("gagal menghapus data mata pelajaran: %w", err)
	}
	// Jika proses hapus berhasil maka kembalikan nil.
	return nil
}
