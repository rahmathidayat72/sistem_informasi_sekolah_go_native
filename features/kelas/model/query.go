package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/kelas"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// kelasQuery merepresentasikan query yang berhubungan dengan tabel kelas.
// Struct ini memiliki satu field yaitu db yang berisi koneksi database.
// Fungsi ini digunakan untuk menghandle query ke database.
type kelasQuery struct {
	// db berisi koneksi database yang digunakan untuk menghandle query.
	db *pgxpool.Pool
}

// NewDataKelas membuat objek kelasQuery yang berisi koneksi database.
// Fungsi ini digunakan untuk menginisialisasi objek kelasQuery yang berisi koneksi database.
// Jika parameter db nil maka akan terjadi panic.
func NewDataKelas(db *pgxpool.Pool) kelas.DataKelasInterface {
	// Jika db nil maka akan terjadi panic
	if db == nil {
		panic("Nil database")
	}
	// Membuat objek kelasQuery baru dan menginisialisasi field db dengan parameter db
	return &kelasQuery{db: db}
}

// Insert implements kelas.DataKelasInterface.
// Fungsi ini digunakan untuk menginsert data kelas ke dalam database.
// Fungsi ini mengembalikan error jika terjadi kesalahan.
func (k *kelasQuery) Insert(insert *kelas.KelasCore) error {
	// Jika parameter db nil maka akan terjadi panic.
	if k.db == nil {
		return errors.New("kelas model: Nil database connection")
	}

	// Jika parameter insert nil maka akan terjadi error.
	if insert == nil {
		return errors.New("insert data is nil")
	}

	// Generate UUID untuk ID kelas jika belum ada.
	if insert.ID == "" {
		insert.ID = uuid.New().String()
	}

	// --- Validasi dan sinkronisasi Nama_Guru & ID_Guru ---
	// Validasi apakah Nama_Guru dan ID_Guru kosong atau tidak.
	// Jika keduanya kosong maka akan terjadi error.
	// Jika hanya Nama_Guru diisi maka cari ID-nya berdasarkan nama guru.
	// Jika hanya ID_Guru diisi maka cari nama-nya berdasarkan ID guru.
	// Jika keduanya diisi maka validasi apakah cocok atau tidak.
	switch {
	case insert.Nama_Guru != "" && insert.ID_Guru == "":
		// Jika hanya Nama_Guru diisi → cari ID-nya
		var idGuru string
		err := k.db.QueryRow(context.Background(), "SELECT id FROM guru WHERE nama = $1", insert.Nama_Guru).Scan(&idGuru)
		if err != nil {
			log.Printf("InsertKelas: nama guru '%s' tidak ditemukan", insert.Nama_Guru)
			return fmt.Errorf("guru dengan nama '%s' tidak ditemukan", insert.Nama_Guru)
		}
		insert.ID_Guru = idGuru

	case insert.ID_Guru != "" && insert.Nama_Guru == "":
		// Jika hanya ID_Guru diisi → cari nama-nya
		var namaGuru string
		err := k.db.QueryRow(context.Background(), "SELECT nama FROM guru WHERE id = $1", insert.ID_Guru).Scan(&namaGuru)
		if err != nil {
			log.Printf("InsertKelas: ID guru '%s' tidak ditemukan", insert.ID_Guru)
			return fmt.Errorf("guru dengan ID '%s' tidak ditemukan", insert.ID_Guru)
		}
		insert.Nama_Guru = namaGuru

	case insert.ID_Guru != "" && insert.Nama_Guru != "":
		// Jika keduanya diisi → validasi apakah cocok
		var existingName string
		err := k.db.QueryRow(context.Background(), "SELECT nama FROM guru WHERE id = $1", insert.ID_Guru).Scan(&existingName)
		if err != nil {
			log.Printf("InsertKelas: ID guru '%s' tidak ditemukan", insert.ID_Guru)
			return fmt.Errorf("guru dengan ID '%s' tidak ditemukan", insert.ID_Guru)
		}
		if strings.TrimSpace(existingName) != strings.TrimSpace(insert.Nama_Guru) {
			log.Printf("InsertKelas: Nama guru tidak cocok dengan ID guru. Dapat: '%s', seharusnya: '%s'", insert.Nama_Guru, existingName)
			return fmt.Errorf("nama guru '%s' tidak cocok dengan ID guru '%s'", insert.Nama_Guru, insert.ID_Guru)
		}
	}

	// --- Siapkan ID_Guru untuk query INSERT (boleh null) ---
	var idGuruParam interface{}
	if insert.ID_Guru == "" {
		idGuruParam = nil
	} else {
		idGuruParam = insert.ID_Guru
	}

	// --- Jalankan query INSERT ---
	query := `INSERT INTO kelas (id, kelas, id_guru) VALUES ($1, $2, $3)`
	_, err := k.db.Exec(context.Background(), query,
		insert.ID,
		insert.Kelas,
		idGuruParam,
	)
	if err != nil {
		log.Printf("InsertKelas error exec: %v", err)
		return fmt.Errorf("insert failed: %w", err)
	}

	return nil
}

// SelectAll digunakan untuk mengambil semua data kelas dari database.
// Fungsi ini mengembalikan slice kelas.KelasCore yang berisi data kelas.
// Jika terjadi error maka fungsi ini akan mengembalikan error.
func (k *kelasQuery) SelectAll() ([]kelas.KelasCore, error) {
	// Validasi koneksi database
	if k.db == nil {
		// Jika koneksi database nil, kembalikan error
		return nil, errors.New("Nil database connection")
	}

	// Query untuk mengambil semua data kelas dan nama guru yang terkait
	query := `
		SELECT 
			k.id, k.kelas, k.id_guru, g.nama AS nama_guru
		FROM 
			kelas k
		LEFT JOIN 
			guru g ON k.id_guru = g.id
		WHERE 
			k.delete_at IS NULL
	`

	// Jalankan query dan simpan hasilnya dalam rows
	rows, err := k.db.Query(context.Background(), query)
	if err != nil {
		// Jika terjadi error saat eksekusi query, log error dan kembalikan
		log.Printf("SelectAll error exec: %v", err)
		return nil, fmt.Errorf("select failed: %w", err)
	}
	defer rows.Close() // Pastikan rows ditutup setelah selesai digunakan

	var result []kelas.KelasCore // Variabel untuk menyimpan hasil

	// Iterasi melalui hasil rows
	// Fungsi ini digunakan untuk mengiterasi data yang diambil dari database dan
	// memasukkan data tersebut dalam slice kelas.KelasCore
	for rows.Next() {
		var kelas Kelas // Deklarasi variabel kelas untuk menyimpan data yang diiterasi

		var idGuru sql.NullString   // Deklarasi variabel idGuru yang dapat bernilai null
		var namaGuru sql.NullString // Deklarasi variabel namaGuru yang dapat bernilai null

		// Pindai setiap baris ke dalam variabel kelas
		// Fungsi Scan digunakan untuk memindai setiap baris yang diiterasi
		// dan menyimpannya dalam variabel kelas
		err := rows.Scan(&kelas.ID, &kelas.Kelas, &idGuru, &namaGuru)
		if err != nil {
			// Jika terjadi error saat scan, log error dan kembalikan
			// Fungsi log.Printf digunakan untuk mencatat log error
			// dan mengembalikan error
			log.Printf("SelectAll error scan: %v", err)
			return nil, fmt.Errorf("select failed: %w", err)
		}

		// Konversi nilai null dari database
		// Fungsi ini digunakan untuk mengkonversi nilai null yang diambil dari database
		// menjadi nilai yang sesuai untuk KelasCore
		if idGuru.Valid {
			kelas.ID_Guru = idGuru.String // Jika idGuru valid maka simpan nilainya dalam kelas.ID_Guru
		}
		if namaGuru.Valid {
			kelas.Nama_Guru = namaGuru.String // Jika namaGuru valid maka simpan nilainya dalam kelas.Nama_Guru
		}

		// Format hasil scan menjadi KelasCore dan simpan dalam result
		// Fungsi FormatterResponse digunakan untuk memformat data yang diiterasi
		// menjadi KelasCore dan menyimpannya dalam result
		core := FormatterResponse(kelas)
		result = append(result, core) // Tambahkan KelasCore yang dihasilkan ke dalam result
	}

	// Cek error setelah iterasi
	if err := rows.Err(); err != nil {
		// Jika terjadi error pada rows, log error dan kembalikan
		log.Printf("SelectAll error rows: %v", err)
		return nil, fmt.Errorf("select failed: %w", err)
	}

	// Log jumlah kelas yang berhasil diambil
	log.Printf("Successfully fetched %d kelas from database", len(result))
	// Kembalikan hasil dalam bentuk slice kelas.KelasCore
	return result, nil
}

// SelectById digunakan untuk mengambil data kelas berdasarkan ID yang diberikan.
// Fungsi ini akan mengembalikan objek kelas.KelasCore yang sesuai dengan ID tersebut
// dan error jika terjadi kesalahan dalam pengambilan data.
func (k *kelasQuery) SelectById(id string) (*kelas.KelasCore, error) {
	// Memeriksa apakah koneksi ke database ada atau tidak
	if k.db == nil {
		// Jika koneksi database nil, kembalikan error
		return nil, errors.New("Nil database connection")
	}

	// Query SQL untuk mengambil data kelas berdasarkan ID dan memastikan data belum dihapus
	query := `SELECT 
			k.id, k.kelas, k.id_guru, g.nama AS nama_guru
		FROM 
			kelas k
		LEFT JOIN 
			guru g ON k.id_guru = g.id 
		WHERE 
			k.id = $1 AND k.delete_at IS NULL`

	var kelas kelas.KelasCore // Deklarasi variabel kelas untuk menyimpan hasil query

	// Eksekusi query dan scan hasilnya ke dalam variabel kelas
	err := k.db.QueryRow(context.Background(), query, id).Scan(
		&kelas.ID,        // Scan kolom id ke dalam kelas.ID
		&kelas.Kelas,     // Scan kolom kelas ke dalam kelas.Kelas
		&kelas.ID_Guru,   // Scan kolom id_guru ke dalam kelas.ID_Guru
		&kelas.Nama_Guru, // Scan kolom nama_guru ke dalam kelas.Nama_Guru
	)
	if err != nil {
		// Jika terjadi error saat eksekusi query, log error dan kembalikan
		log.Printf("SelectById error exec: %v", err)
		return nil, fmt.Errorf("select failed: %w", err)
	}

	// Jika berhasil, kembalikan data kelas
	return &kelas, nil
}

// Update digunakan untuk mengupdate data kelas berdasarkan ID yang diberikan.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan dalam proses update.
func (k *kelasQuery) Update(insert *kelas.KelasCore, id string) error {
	// Memeriksa apakah koneksi ke database ada atau tidak
	if k.db == nil {
		// Jika koneksi database nil, kembalikan error
		return errors.New("Nil database connection")
	}

	// Query SQL untuk mengupdate data kelas berdasarkan ID
	query := `UPDATE kelas SET kelas = $1, id_guru = $2 WHERE id = $3`

	// Eksekusi query update dengan parameter yang diberikan
	res, err := k.db.Exec(context.Background(), query,
		insert.Kelas,   // Menggunakan nilai kelas baru dari parameter insert
		insert.ID_Guru, // Menggunakan ID guru baru dari parameter insert
		id,             // ID dari kelas yang akan diupdate
	)
	if err != nil {
		// Jika terjadi error saat eksekusi query, log error dan kembalikan
		log.Printf("UpdateKelas error exec: %v", err)
		return fmt.Errorf("update failed: %w", err)
	}

	// Memeriksa apakah ada baris yang terpengaruh oleh update
	if res.RowsAffected() == 0 {
		// Jika tidak ada baris yang terpengaruh, log dan kembalikan error
		log.Printf("Updatekelas: no rows updated for id %s", id)
		return errors.New("update failed: no rows affected")
	}

	// Kembalikan nil jika update berhasil tanpa error
	return nil
}

// DeleteById implements kelas.DataKelasInterface.
// DeleteById implements kelas.DataKelasInterface.
// Fungsi ini digunakan untuk menghapus data kelas berdasarkan ID yang diberikan.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan dalam proses hapus.
func (k *kelasQuery) DeleteById(id string) error {
	// Memeriksa apakah koneksi ke database ada atau tidak
	if k.db == nil {
		// Jika koneksi database nil, kembalikan error
		return errors.New("Nil database connection")
	}

	// Query SQL untuk menghapus data kelas berdasarkan ID
	// dengan menggunakan soft delete, yaitu mengupdate kolom delete_at menjadi NOW()
	query := "UPDATE kelas SET delete_at = NOW() WHERE id = $1 AND delete_at IS NULL"

	// Eksekusi query delete dengan parameter yang diberikan
	res, err := k.db.Exec(context.Background(), query, id)
	if err != nil {
		// Jika terjadi error saat eksekusi query, log error dan kembalikan
		log.Printf("DeleteById error exec: %v", err)
		return fmt.Errorf("delete failed: %w", err)
	}

	// Memeriksa apakah ada baris yang terpengaruh oleh delete
	if res.RowsAffected() == 0 {
		// Jika tidak ada baris yang terpengaruh, log dan kembalikan error
		log.Printf("DeleteById: no rows deleted for id %s", id)
		return errors.New("delete failed: no rows affected")
	}

	// Kembalikan nil jika delete berhasil tanpa error
	return nil
}
