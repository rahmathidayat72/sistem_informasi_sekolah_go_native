package model

import (
	"context"
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/users"
	"go_rest_native_sekolah/helper"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserQuerry adalah struktur data yang berisi pointer ke objek pgxpool.Pool.
// Struktur data ini digunakan untuk menghandle query ke database yang berhubungan dengan tabel users.
type UserQuerry struct {
	db *pgxpool.Pool // db adalah pointer ke objek pgxpool.Pool yang digunakan untuk menghandle query ke database.
}

// NewUserData membuat objek UserQuerry yang berisi koneksi database.
// Fungsi ini digunakan untuk menginisialisasi objek UserQuerry yang berisi koneksi database.
// Jika parameter db nil maka akan terjadi panic.
func NewUserData(db *pgxpool.Pool) users.DataUserInterface {
	// Jika db nil maka akan terjadi panic
	if db == nil {
		panic("Nil database")
	}
	// Membuat objek UserQuerry baru dan menginisialisasi field db dengan parameter db
	return &UserQuerry{db: db}
}

// SelectAllUser mengembalikan slice users.UserCore yang berisi semua data user
// yang tersimpan di database.
// Fungsi ini mengembalikan error jika terjadi kesalahan saat query ke database.
func (u *UserQuerry) SelectAllUser() ([]users.UserCore, error) {
	if u.db == nil {
		// Jika koneksi database tidak ada maka kembalikan error
		return nil, errors.New("Koneksi database tidak ada")
	}

	// Query untuk mengambil semua data user dari database
	// Filter data user yang tidak dihapus
	query := "SELECT id, username, email, password, role FROM users WHERE delete_at IS NULL"

	rows, err := u.db.Query(context.Background(), query)
	if err != nil {
		// Jika terjadi error saat query maka kembalikan error
		return nil, err
	}
	defer rows.Close() // Pastikan rows ditutup setelah selesai digunakan

	var result []users.UserCore // Variabel untuk menyimpan hasil

	// Iterasi melalui hasil rows
	// Fungsi ini digunakan untuk mengiterasi data yang diambil dari database dan
	// memasukkan data tersebut dalam slice users.UserCore
	for rows.Next() {
		var user users.UserCore // Deklarasi variabel user untuk menyimpan data yang diiterasi

		// Pindai setiap baris ke dalam variabel user
		// Fungsi Scan digunakan untuk memindai setiap baris yang diiterasi
		// dan menyimpannya dalam variabel user
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
		if err != nil {
			// Jika terjadi error saat scan maka kembalikan error
			return nil, err
		}

		// Format data user yang diiterasi menjadi objek UserCore yang sesuai
		core := FormatterRequest(user)
		result = append(result, FormatterResponse(core))
	}

	// Cek apakah ada error setelah loop
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Cek apakah ada data yang diiterasi
	if len(result) == 0 {
		// Jika tidak ada data maka kembalikan error
		return nil, errors.New("Tidak ada data user yang ditemukan di database")
	}

	// Log suksesnya query dan kembalikan hasil
	log.Printf("Berhasil mengambil %d users dari database", len(result))

	return result, nil
}

// InsertUser implements users.DataUserInterface.
// Fungsi ini digunakan untuk menginsert data user ke dalam database.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (u *UserQuerry) InsertUser(input *users.UserCore) error {
	// Cek apakah objek UserQuerry atau koneksi database adalah nil.
	if u == nil || u.db == nil {
		return errors.New("Nil UserQuerry or database")
	}
	// Cek apakah input user adalah nil.
	if input == nil {
		return errors.New("Nil input")
	}

	// Generate UUID jika ID belum ada.
	if input.ID == "" {
		input.ID = uuid.New().String()
	}

	// Mapping data user ke struct User untuk keperluan penyimpanan di database.
	userInput := FormatterRequest(*input)

	// Hash password dari input.Password agar tersimpan dengan aman di database.
	hashedPassword := helper.HashPassword(input.Password)
	userInput.Password = hashedPassword

	// Query untuk memasukkan data user ke dalam tabel users.
	query := `INSERT INTO users (id, username, email, password, role) VALUES ($1, $2, $3, $4, $5)`

	// Eksekusi query untuk menyimpan data user ke dalam database.
	_, err := u.db.Exec(context.Background(), query,
		userInput.ID,
		userInput.Username,
		userInput.Email,
		userInput.Password,
		userInput.Role,
	)
	// Jika terjadi error saat eksekusi query, log error dan kembalikan sebagai hasil fungsi.
	if err != nil {
		log.Printf("InsertUser error exec: %v", err)
		return fmt.Errorf("insert failed: %w", err)
	}

	// Kembalikan nil jika insert berhasil tanpa error.
	return nil
}

// SelectUserById implements users.DataUserInterface.
// Fungsi ini digunakan untuk mengambil data user berdasarkan id
// yang dikirimkan sebagai parameter.
// Fungsi ini akan mengembalikan pointer ke struct UserCore yang berisi data user
// dan error jika terjadi kesalahan.
func (u *UserQuerry) SelectUserById(id string) (*users.UserCore, error) {
	// Cek apakah objek UserQuerry atau koneksi database adalah nil.
	// Jika nil maka kembalikan error.
	if u == nil || u.db == nil {
		return nil, errors.New("Nil UserQuerry or database")
	}

	// Cek apakah id yang dikirimkan adalah kosong.
	// Jika kosong maka kembalikan error.
	if id == "" {
		return nil, errors.New("Nil or empty id")
	}

	// Query untuk mengambil data user berdasarkan id.
	// Query ini akan mengambil kolom id, username, email, password, dan role
	// berdasarkan id yang dikirimkan dan delete_at IS NULL
	// yang artinya data user yang diambil belum dihapus.
	query := "SELECT id, username, email, password, role FROM users WHERE id = $1 AND delete_at IS NULL"

	// Jalankan query.
	// Fungsi QueryRow akan mengembalikan row yang sesuai dengan query
	// dan error jika terjadi kesalahan.
	row := u.db.QueryRow(context.Background(), query, id)

	// Deklarasikan variabel result yang akan digunakan untuk menyimpan hasil query.
	var result users.UserCore

	// Scan hasil query ke variabel result.
	// Fungsi Scan akan mengembalikan error jika terjadi kesalahan.
	err := row.Scan(&result.ID, &result.Username, &result.Email, &result.Password, &result.Role)
	if err != nil {
		// Jika terjadi error maka kembalikan error.
		return nil, fmt.Errorf("gagal mengambil data user: %w", err)
	}

	// Kembalikan data user yang diambil.
	return &result, nil
}

// UpdateUser implements users.DataUserInterface.
// Fungsi ini digunakan untuk mengupdate data user berdasarkan ID yang diberikan.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan saat proses update.
func (u *UserQuerry) UpdateUser(insert *users.UserCore, id string) error {
	// Memeriksa apakah objek UserQuerry atau koneksi database adalah nil.
	if u == nil || u.db == nil {
		return errors.New("Nil UserQuerry or database")
	}
	// Memeriksa apakah input user adalah nil.
	if insert == nil {
		return errors.New("Nil user")
	}
	// Memeriksa apakah ID yang diberikan kosong.
	if id == "" {
		return errors.New("Nil or empty id")
	}
	// Hash password sebelum simpan ke database untuk keamanan.
	hashedPassword := helper.HashPassword(insert.Password)

	// Membuat query SQL untuk mengupdate data user berdasarkan ID.
	query := "UPDATE users SET username = $2, email = $3, password = $4, role = $5 WHERE id = $1"
	// Menjalankan query update pada database dengan parameter yang diberikan.
	res, err := u.db.Exec(context.Background(), query, id, insert.Username, insert.Email, hashedPassword, insert.Role)
	if err != nil {
		// Log error jika terjadi kesalahan saat eksekusi query.
		log.Printf("UpdateUser error exec: %v", err)
		return fmt.Errorf("update failed: %w", err)
	}

	// Memeriksa apakah ada baris yang terpengaruh oleh update.
	if res.RowsAffected() == 0 {
		// Log dan kembalikan error jika tidak ada baris yang terpengaruh.
		log.Printf("UpdateUser: no rows updated for id %s", id)
		return errors.New("update failed: no rows affected")
	}

	// Mengembalikan nil jika update berhasil tanpa error.
	return nil
}

// DeleteUserById implements users.DataUserInterface.
// DeleteUserById implements users.DataUserInterface.
// Fungsi ini digunakan untuk menghapus data user berdasarkan ID yang diberikan.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan saat proses hapus.
func (u *UserQuerry) DeleteUserById(id string) error {
	// Memeriksa apakah objek UserQuerry atau koneksi database adalah nil.
	// Jika nil maka kembalikan error.
	if u == nil || u.db == nil {
		return errors.New("Nil UserQuerry or database")
	}
	// Memeriksa apakah ID yang diberikan kosong.
	// Jika kosong maka kembalikan error.
	if id == "" {
		return errors.New("Nil or empty id")
	}

	// Membuat query SQL untuk mengupdate user berdasarkan ID.
	// Query ini menggunakan soft delete, yaitu mengupdate kolom delete_at menjadi NOW()
	// jika data user dengan ID yang dikirimkan memang ada dan belum dihapus.
	query := "UPDATE users SET delete_at = NOW() WHERE id = $1"

	// Menjalankan query update pada database dengan parameter yang diberikan.
	res, err := u.db.Exec(context.Background(), query, id)
	if err != nil {
		// Log error jika terjadi kesalahan saat eksekusi query.
		log.Printf("DeleteUserById error exec: %v", err)
		return fmt.Errorf("delete failed: %w", err)
	}

	// Memeriksa apakah ada baris yang terpengaruh oleh update.
	// Jika tidak ada baris yang terpengaruh maka log dan kembalikan error.
	if res.RowsAffected() == 0 {
		log.Printf("DeleteUserById: no rows updated for id %s", id)
		return errors.New("delete failed: no rows affected")
	}

	// Mengembalikan nil jika hapus berhasil tanpa error.
	return nil
}
