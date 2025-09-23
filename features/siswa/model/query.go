package model

import (
	"context"
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/siswa"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// siswaQuery adalah struct yang digunakan untuk menghandle query ke database yang berhubungan dengan tabel siswa.
// Struct ini memiliki satu field yaitu db yang berisi koneksi database.
type siswaQuery struct {
	db *pgxpool.Pool // Koneksi database yang digunakan untuk menghandle query ke database.
}

// NewSiswaData membuat objek siswaQuery yang berisi koneksi database.
// Fungsi ini digunakan untuk menginisialisasi objek siswaQuery yang berisi koneksi database.
// Jika parameter db nil maka akan terjadi panic.
func NewSiswaData(db *pgxpool.Pool) siswa.DataSiswaInterface {
	// Jika db nil maka akan terjadi panic
	if db == nil {
		panic("Nil database")
	}
	// Membuat objek siswaQuery baru dan menginisialisasi field db dengan parameter db
	return &siswaQuery{db: db}
}

// InsertSiswa implements siswa.DataSiswaInterface.
// InsertSiswa adalah fungsi yang digunakan untuk menginsert data siswa ke dalam database.
// Fungsi ini menerima parameter objek siswa.SiswaCore yang berisi data-data siswa yang akan diinsert.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (s *siswaQuery) InsertSiswa(insert *siswa.SiswaCore) error {
	if s.db == nil {
		// Jika parameter db nil maka akan terjadi panic.
		return errors.New("Nil database")
	}

	if insert == nil {
		// Jika parameter insert nil maka akan terjadi error.
		return errors.New("insert data is nil")
	}

	if insert.ID == "" {
		// Jika ID belum diisi, generate UUID baru.
		insert.ID = uuid.New().String()
	}

	// --- Validasi dan sinkronisasi Nama_Kelas & Kelas_ID ---
	// Validasi apakah Nama_Kelas dan Kelas_ID kosong atau tidak.
	// Jika keduanya kosong maka akan terjadi error.
	// Jika hanya Nama_Kelas diisi maka cari ID-nya berdasarkan nama kelas.
	// Jika hanya Kelas_ID diisi maka cari nama-nya berdasarkan ID kelas.
	// Jika keduanya diisi maka validasi apakah cocok atau tidak.
	switch {
	case insert.Nama_Kelas != "" && insert.Kelas_ID == "":
		// Jika hanya Nama_Kelas diisi → cari Kelas_ID-nya
		// Trim nama kelas agar tidak ada spasi di awal dan akhir.
		insert.Nama_Kelas = strings.TrimSpace(insert.Nama_Kelas)

		var kelasID string
		err := s.db.QueryRow(context.Background(),
			"SELECT id FROM kelas WHERE TRIM(kelas) ILIKE TRIM($1)", insert.Nama_Kelas).Scan(&kelasID)
		if err != nil {
			// Jika tidak ada kelas dengan nama yang sesuai maka akan terjadi error.
			log.Printf("InsertKelas: nama kelas '%s' tidak ditemukan", insert.Nama_Kelas)
			return fmt.Errorf("kelas dengan nama '%s' tidak ditemukan", insert.Nama_Kelas)
		}
		// Masukkan ID kelas ke dalam objek siswa.
		insert.Kelas_ID = kelasID

	case insert.Kelas_ID != "" && insert.Nama_Kelas == "":
		// Jika hanya Kelas_ID diisi → cari nama-nya
		var namaKelas string
		err := s.db.QueryRow(context.Background(),
			"SELECT kelas FROM kelas WHERE id = $1", insert.Kelas_ID).Scan(&namaKelas)
		if err != nil {
			// Jika tidak ada kelas dengan ID yang sesuai maka akan terjadi error.
			log.Printf("InsertKelas: ID kelas '%s' tidak ditemukan", insert.Kelas_ID)
			return fmt.Errorf("kelas dengan ID '%s' tidak ditemukan", insert.Kelas_ID)
		}
		// Masukkan nama kelas ke dalam objek siswa.
		insert.Nama_Kelas = namaKelas

	case insert.Kelas_ID != "" && insert.Nama_Kelas != "":
		// Jika keduanya diisi → validasi apakah cocok
		// Trim nama kelas agar tidak ada spasi di awal dan akhir.
		insert.Nama_Kelas = strings.TrimSpace(insert.Nama_Kelas)

		var existingName string
		err := s.db.QueryRow(context.Background(),
			"SELECT kelas FROM kelas WHERE id = $1", insert.Kelas_ID).Scan(&existingName)
		if err != nil {
			// Jika tidak ada kelas dengan ID yang sesuai maka akan terjadi error.
			log.Printf("InsertKelas: ID kelas '%s' tidak ditemukan", insert.Kelas_ID)
			return fmt.Errorf("kelas dengan ID '%s' tidak ditemukan", insert.Kelas_ID)
		}
		// Validasi apakah nama kelas yang diinput sama dengan nama kelas yang ada di database.
		if strings.TrimSpace(strings.ToLower(existingName)) != strings.ToLower(insert.Nama_Kelas) {
			// Jika tidak sama maka akan terjadi error.
			log.Printf("InsertKelas: Nama kelas tidak cocok dengan ID kelas. Dapat: '%s', seharusnya: '%s'",
				insert.Nama_Kelas, existingName)
			return fmt.Errorf("nama kelas '%s' tidak cocok dengan ID kelas '%s'", insert.Nama_Kelas, insert.Kelas_ID)
		}
	}

	// --- Siapkan Kelas_ID untuk query INSERT (boleh null) ---
	// Jika Kelas_ID kosong maka akan diisi dengan nilai null.
	var idKelasParam interface{}
	if insert.Kelas_ID == "" {
		idKelasParam = nil
	} else {
		idKelasParam = insert.Kelas_ID
	}

	// --- Eksekusi query INSERT ke tabel siswa ---
	_, err := s.db.Exec(context.Background(),
		"INSERT INTO siswa (id, kelas_id, nama, email, alamat) VALUES ($1, $2, $3, $4, $5)",
		insert.ID, idKelasParam, insert.Nama, insert.Email, insert.Alamat)
	if err != nil {
		// Jika terjadi kesalahan maka akan terjadi error.
		log.Printf("InsertSiswa error exec: %v", err)
		return fmt.Errorf("insert failed: %w", err)
	}

	return nil
}

// SelectAllSiswa implements siswa.DataSiswaInterface.
// Fungsi ini digunakan untuk mengambil seluruh data siswa dari database.
// Fungsi ini akan mengembalikan array siswa.SiswaCore yang berisi data-data siswa.
// Jika terjadi kesalahan maka akan mengembalikan error.
func (s *siswaQuery) SelectAllSiswa() ([]siswa.SiswaCore, error) {
	if s.db == nil {
		// Jika koneksi database tidak ada maka akan mengembalikan error.
		return nil, errors.New("Nil database")
	}

	query := `SELECT 
    s.id, 
    s.kelas_id, 
    k.kelas AS nama_kelas,
    s.nama, 
    s.email, 
    s.alamat
FROM 
    siswa s
LEFT JOIN 
    kelas k ON s.kelas_id = k.id
WHERE 
    s.delete_at IS NULL`

	// Eksekusi query ke database.
	rows, err := s.db.Query(context.Background(), query)
	if err != nil {
		// Jika terjadi kesalahan maka akan mengembalikan error.
		log.Printf("SelectAllSiswa error query: %v", err)
		return nil, fmt.Errorf("select failed: %w", err)
	}

	// Tutup koneksi database setelah selesai digunakan.
	defer rows.Close()

	// Deklarasikan variabel result yang akan digunakan untuk menyimpan hasil query.
	var result []siswa.SiswaCore

	// Looping untuk mengambil data siswa dari hasil query.
	// Setiap data siswa akan di ambil dan di format ke dalam objek siswa.SiswaCore.
	for rows.Next() {
		// Deklarasikan variabel siswa yang akan digunakan untuk menyimpan data siswa.
		var siswa Siswa

		// Ambil data siswa dari hasil query dan simpan ke dalam objek siswa.
		err := rows.Scan(&siswa.ID, &siswa.Kelas_ID, &siswa.Nama_Kelas, &siswa.Nama, &siswa.Email, &siswa.Alamat)
		if err != nil {
			// Jika terjadi kesalahan maka akan mengembalikan error.
			log.Printf("SelectAllSiswa error scan: %v", err)
			return nil, fmt.Errorf("select failed: %w", err)
		}

		// Format data siswa ke dalam objek siswa.SiswaCore.
		core := FormatterResponse(siswa)

		// Tambahkan data siswa ke dalam array result.
		result = append(result, core)
	}

	// Jika terjadi kesalahan maka akan mengembalikan error.
	if err := rows.Err(); err != nil {
		log.Printf("SelectAll error rows: %v", err)
		return nil, fmt.Errorf("select failed: %w", err)
	}

	// Log berapa banyak data siswa yang berhasil diambil.
	log.Printf("Successfully fetched %d kelas from database", len(result))

	// Mengembalikan array result yang berisi data-data siswa.
	return result, nil
}

// SelectById implements siswa.DataSiswaInterface.
// SelectById implements siswa.DataSiswaInterface.
// Fungsi ini digunakan untuk mengambil data siswa berdasarkan ID.
// Fungsi ini akan mengembalikan data siswa yang sesuai dengan ID yang dikirimkan
// dan error jika terjadi kesalahan.
func (s *siswaQuery) SelectById(id string) (*siswa.SiswaCore, error) {
	if s == nil || s.db == nil {
		// Jika koneksi database tidak ada maka kembalikan error.
		return nil, errors.New("Nil database")
	}
	if id == "" {
		// Jika ID kosong maka kembalikan error.
		return nil, errors.New("ID cannot be empty")
	}

	// Query untuk mengambil data siswa berdasarkan ID.
	// Query ini akan mengambil kolom id, kelas_id, nama_kelas, nama, email, dan alamat
	// berdasarkan ID yang dikirimkan dan delete_at IS NULL
	// yang artinya data siswa yang diambil belum dihapus.
	query := "SELECT s.id, s.kelas_id, k.kelas AS nama_kelas, s.nama, s.email, s.alamat FROM siswa s LEFT JOIN kelas k ON s.kelas_id = k.id WHERE s.id = $1 AND s.delete_at IS NULL"

	// Jalankan query.
	// Fungsi QueryRow akan mengembalikan row yang sesuai dengan query
	// dan error jika terjadi kesalahan.
	row := s.db.QueryRow(context.Background(), query, id)

	// Deklarasikan variabel result yang akan digunakan untuk menyimpan hasil query.
	var result siswa.SiswaCore

	// Scan hasil query ke variabel result.
	// Fungsi Scan akan mengembalikan error jika terjadi kesalahan.
	err := row.Scan(&result.ID, &result.Kelas_ID, &result.Nama_Kelas, &result.Nama, &result.Email, &result.Alamat)
	if err != nil {
		// Jika terjadi kesalahan maka kembalikan error.
		log.Printf("SelectById error scan: %v", err)
		return nil, fmt.Errorf("select failed: %w", err)
	}

	// Log berapa banyak data siswa yang berhasil diambil.
	log.Printf("Successfully fetched siswa from database with id:%s", id)

	// Mengembalikan data siswa yang diambil.
	return &result, nil
}

// Update implements siswa.DataSiswaInterface.
// Update implements siswa.DataSiswaInterface.
// Fungsi ini digunakan untuk mengupdate data siswa berdasarkan ID.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (s *siswaQuery) Update(insert *siswa.SiswaCore, id string) error {
	// Cek apakah koneksi database ada atau tidak.
	if s == nil || s.db == nil {
		return errors.New("Nil database")
	}
	// Cek apakah ID kosong atau tidak.
	if id == "" {
		return errors.New("ID tidak boleh kosong")
	}

	// Query untuk mengupdate data siswa berdasarkan ID.
	// Query ini akan mengupdate kolom nama, email, alamat, dan kelas_id.
	// berdasarkan ID yang dikirimkan dan delete_at IS NULL
	// yang artinya data siswa yang diupdate belum dihapus.
	query := "UPDATE siswa SET nama = $1, email = $2, alamat = $3, kelas_id = $4 WHERE id = $5"
	// Jalankan query untuk mengupdate data siswa.
	// Fungsi Exec digunakan untuk mengeksekusi query yang tidak mengembalikan hasil.
	res, err := s.db.Exec(context.Background(), query, insert.Nama, insert.Email, insert.Alamat, insert.Kelas_ID, id)
	if err != nil {
		// Jika terjadi error saat query maka log error dan kembalikan.
		log.Printf("Update error exec: %v", err)
		return fmt.Errorf("update failed: %w", err)
	}
	// Cek apakah ada baris yang terpengaruh.
	if res.RowsAffected() == 0 {
		// Jika tidak ada baris yang terpengaruh maka log dan kembalikan error.
		log.Printf("UpdateUser: no rows updated for id %s", id)
		return errors.New("update failed: no rows affected")
	}
	// Log berapa banyak data siswa yang berhasil diupdate.
	log.Printf("Successfully updated %s in database", id)
	// Mengembalikan nil jika update berhasil tanpa error.
	return nil
}

// DeleteById implements siswa.DataSiswaInterface.
// DeleteById implements siswa.DataSiswaInterface.
// Fungsi ini digunakan untuk menghapus data siswa berdasarkan ID.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (s *siswaQuery) DeleteById(id string) error {
	// Cek apakah koneksi database ada atau tidak.
	// Jika tidak ada maka kembalikan error.
	if s.db == nil {
		return errors.New("Koneksi database tidak ada")
	}

	// Query untuk menghapus data siswa berdasarkan ID.
	// Query ini menggunakan parameter $1 untuk menggantikan nilai id.
	// Query ini akan mengupdate kolom delete_at dengan waktu sekarang
	// jika data siswa dengan ID yang dikirimkan memang ada dan belum dihapus.
	query := "UPDATE siswa SET delete_at = NOW() WHERE id = $1 AND delete_at IS NULL"

	// Jalankan query untuk menghapus data siswa.
	// Fungsi Exec digunakan untuk mengeksekusi query yang tidak mengembalikan hasil.
	// Fungsi Exec juga akan mengembalikan error jika terjadi kesalahan.
	res, err := s.db.Exec(context.Background(), query, id)
	if err != nil {
		// Jika terjadi error saat query maka log error dan kembalikan.
		log.Printf("DeleteById error exec: %v", err)
		return fmt.Errorf("hapus gagal: %w", err)
	}

	// Cek apakah ada baris yang terpengaruh.
	// Jika tidak ada baris yang terpengaruh maka log dan kembalikan error.
	if res.RowsAffected() == 0 {
		log.Printf("DeleteById: tidak ada baris yang dihapus untuk id %s", id)
		return errors.New("hapus gagal: tidak ada baris yang terpengaruh")
	}

	// Jika data berhasil dihapus maka log pesan sukses dan kembalikan nil.
	log.Printf("Berhasil menghapus siswa dengan id: %s", id)
	return nil
}
