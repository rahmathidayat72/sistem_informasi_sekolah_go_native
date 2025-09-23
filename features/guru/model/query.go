package gurumodels

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/guru"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// guruQuery merepresentasikan query yang berhubungan dengan tabel guru.
// guruQuery menggunakan database PostgreSQL untuk menghandle query ke database.
type guruQuery struct {
	db *pgxpool.Pool
}

// NewDataGuru membuat objek guruQuery dengan parameter db.
// guruQuery digunakan untuk menghandle query ke database yang berhubungan dengan tabel guru.
// Jika parameter db nil maka akan terjadi panic.
func NewDataGuru(db *pgxpool.Pool) guru.DataGuruInterface {
	if db == nil {
		panic("guru model: Nil database")
	}
	return &guruQuery{db: db}
}

// SelectAllGuru digunakan untuk mengambil semua data guru dari database.
// Fungsi ini akan mengembalikan slice guru.GuruCore yang berisi data guru.
// Jika terjadi error maka fungsi ini akan mengembalikan error.
func (r *guruQuery) SelectAllGuru() ([]guru.GuruCore, error) {
	// Validasi apakah database nil
	if r.db == nil {
		return nil, errors.New("guru model: Nil database")
	}

	// Query untuk mengambil semua data guru
	query := "SELECT id, id_user, nama, email, alamat FROM guru WHERE delete_at IS NULL"

	// Jalankan query
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []guru.GuruCore

	// Looping hasil rows
	for rows.Next() {
		var guru Guru

		// id_user bisa NULL, jadi gunakan sql.NullString
		var idUser sql.NullString

		err := rows.Scan(&guru.ID, &idUser, &guru.Nama, &guru.Email, &guru.Alamat)
		if err != nil {
			log.Printf("SelectAll error scan: %v", err)
			return nil, fmt.Errorf("select failed: %w", err)
		}

		// Konversi idUser ke string biasa
		if idUser.Valid {
			guru.ID_User = idUser.String
		} else {
			guru.ID_User = "" // atau bisa nil jika pakai pointer
		}

		// Format ke Core
		core := FormatterResponse(guru)
		result = append(result, core)
	}

	// Cek error setelah loop
	if err := rows.Err(); err != nil {
		log.Printf("SelectAll error rows: %v", err)
		return nil, fmt.Errorf("select failed: %w", err)
	}

	log.Printf("Successfully fetched %d guru from database", len(result))
	return result, nil
}

// InsertGuru implements guru.DataGuruInterface.
// Fungsi ini digunakan untuk menginsert data guru ke dalam database.
// Fungsi ini mengembalikan error jika terjadi kesalahan.
func (r *guruQuery) InsertGuru(insert *guru.GuruCore) error {
	if r.db == nil {
		return errors.New("guru model: nil database connection")
	}

	// Generate UUID untuk ID guru jika belum ada
	if insert.ID == "" {
		insert.ID = uuid.New().String()
	}

	// Cek apakah ID_User kosong, untuk disisipkan sebagai NULL jika iya
	var idUserParam interface{}
	if insert.ID_User == "" {
		idUserParam = nil
	} else {
		idUserParam = insert.ID_User
	}

	// Query untuk menyimpan data guru
	query := `INSERT INTO guru (id, id_user, nama, email, alamat) VALUES ($1, $2, $3, $4, $5)`

	// Jalankan query
	_, err := r.db.Exec(context.Background(), query,
		insert.ID,
		idUserParam,
		insert.Nama,
		insert.Email,
		insert.Alamat,
	)
	if err != nil {
		log.Printf("InsertGuru error exec: %v", err)
		return fmt.Errorf("insert failed: %w", err)
	}

	return nil
}

// Update implements guru.DataGuruInterface.
// Fungsi ini digunakan untuk mengupdate data guru berdasarkan ID.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (r *guruQuery) Update(insert *guru.GuruCore, id string) error {
	// Cek apakah koneksi database nil
	if r.db == nil {
		return errors.New("koneksi database nil")
	}

	// Cek apakah input data nil
	if insert == nil {
		return errors.New("input data nil")
	}

	// Cek apakah ID kosong
	if id == "" {
		return errors.New("validation error: id harus diisi")
	}

	// Query untuk mengupdate data guru berdasarkan ID
	// query ini akan mengupdate kolom nama, email, dan alamat
	// berdasarkan ID yang dikirimkan
	query := `UPDATE guru SET nama = $2, email = $3, alamat = $4 WHERE id = $1`

	// Eksekusi query update
	// fungsi Exec akan mengembalikan hasil query dan error
	// jika terjadi error maka akan dikembalikan error
	res, err := r.db.Exec(context.Background(), query,
		id,
		insert.Nama,
		insert.Email,
		insert.Alamat,
	)
	if err != nil {
		// Log error jika terjadi kesalahan
		log.Printf("UpdateGuru error exec: %v", err)
		return fmt.Errorf("update failed: %w", err)
	}

	// Cek apakah ada baris yang terpengaruh
	// jika tidak ada baris yang terpengaruh maka akan dikembalikan error
	if res.RowsAffected() == 0 {
		log.Printf("UpdateGuru: no rows updated for id %s", id)
		return errors.New("update failed: no rows affected")
	}

	return nil // Jika tidak ada error maka kembalikan nil
}

// SelectById digunakan untuk mengambil data guru berdasarkan ID
// Fungsi ini akan mengembalikan data guru yang sesuai dengan ID yang dikirimkan
// dan error jika terjadi kesalahan
func (r *guruQuery) SelectById(id string) (*guru.GuruCore, error) {
	// Cek koneksi database
	// Jika koneksi database nil maka kembalikan error
	if r.db == nil {
		return nil, errors.New("guru query: koneksi database nil")
	}

	// Query untuk mengambil data guru berdasarkan ID
	// Query ini akan mengambil kolom id, id_user, nama, email, dan alamat
	// berdasarkan ID yang dikirimkan dan delete_at IS NULL
	// yang artinya data guru yang diambil belum dihapus
	query := `
		SELECT id, id_user, nama, email, alamat 
		FROM guru 
		WHERE id = $1 AND delete_at IS NULL
	`

	// Jalankan query
	// Fungsi QueryRow akan mengembalikan row yang sesuai dengan query
	// dan error jika terjadi kesalahan
	row := r.db.QueryRow(context.Background(), query, id)

	// Deklarasikan variabel untuk menyimpan hasil query
	var result guru.GuruCore

	// Scan hasil query ke variabel result
	// Fungsi Scan akan mengembalikan error jika terjadi kesalahan
	err := row.Scan(&result.ID, &result.ID_User, &result.Nama, &result.Email, &result.Alamat)
	if err != nil {
		// Jika terjadi error maka kembalikan error
		return nil, fmt.Errorf("gagal mengambil data guru: %w", err)
	}

	// Jika tidak ada error maka kembalikan data guru
	return &result, nil
}

// DeleteById implements guru.DataGuruInterface.
func (r *guruQuery) DeleteById(id string) error {
	// Cek koneksi database
	if r.db == nil {
		return errors.New("guru query: koneksi database nil")
	}

	// Query untuk menghapus data guru berdasarkan ID
	query := "UPDATE guru SET delete_at = NOW() WHERE id = $1 AND delete_at IS NULL"

	// Eksekusi query
	res, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		log.Printf("DeleteById error exec: %v", err)
		return fmt.Errorf("delete failed: %w", err)
	}

	// Cek apakah ada baris yang terpengaruh
	if res.RowsAffected() == 0 {
		log.Printf("DeleteById: no rows deleted for id %s", id)
		return errors.New("delete failed: no rows affected")
	}

	return nil // Jika tidak ada error maka kembalikan nil
}
