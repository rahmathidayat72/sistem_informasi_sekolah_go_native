package service

import (
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/users"
	"regexp"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// userService adalah struct yang berisi property userData.
// Property userData memiliki tipe interface DataUserInterface
// yang digunakan untuk mengakses data user dari repository.
type userService struct {
	userData users.DataUserInterface // Interface untuk mengakses data user.
}

// NewServiceUser adalah fungsi yang digunakan untuk
// membuat instance baru dari userService.
// Fungsi ini menerima parameter repo dan db.
func NewServiceUser(repo users.DataUserInterface, db *pgxpool.Pool) users.ServiceUserInterface {
	if repo == nil {
		// Jika parameter repo nil maka akan terjadi panic.
		panic("Nil repository")
	}
	// Mengembalikan pointer ke objek userService yang baru
	// dengan menginisialisasi property userData.
	return &userService{userData: repo}
}

// SelectAllUser mengembalikan slice users.UserCore yang berisi semua data user
// yang tersimpan di database.
// Fungsi ini mengembalikan error jika terjadi kesalahan saat query ke database.
func (u *userService) SelectAllUser() ([]users.UserCore, error) {
	if u == nil || u.userData == nil {
		// Jika objek userService atau repository adalah nil
		// maka kembalikan error.
		return nil, errors.New("Nil service or repository")
	}

	users, err := u.userData.SelectAllUser()
	if err != nil {
		// Jika terjadi error saat query maka kembalikan error
		// dengan menggabungkan pesan error yang diterima.
		return nil, fmt.Errorf("error retrieving users: %w", err)
	}

	// Kembalikan slice users yang diperoleh dari repository.
	return users, nil
}

// InsertUser implements users.ServiceUserInterface.
// Fungsi ini digunakan untuk menginsert data user ke dalam database.
// Fungsi ini menerima parameter insert yang berisi data user yang ingin diinsert.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan saat query ke database.
func (u *userService) InsertUser(insert *users.UserCore) error {
	if u == nil {
		// Jika objek userService adalah nil maka kembalikan error.
		return errors.New("user service: Nil service")
	}
	if u.userData == nil {
		// Jika repository userData adalah nil maka kembalikan error.
		return errors.New("user service: Nil Repository")
	}
	// Periksa apakah input data nil
	if insert == nil {
		// Jika input data adalah nil maka kembalikan error.
		return errors.New("user service: input is nil")
	}

	if insert.Username == "" || insert.Email == "" || insert.Password == "" || insert.Role == "" {
		// Jika salah satu field username, email, password atau role kosong
		// maka kembalikan error.
		return errors.New("validation error: username, email, password dan role harus diisi")
	}
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(insert.Email) {
		// Jika format email tidak sesuai maka kembalikan error.
		return errors.New("validation error: email tidak valid")
	}

	// Panggil fungsi InsertUser pada repository untuk menginsert data user.
	return u.userData.InsertUser(insert)
}

// SelectUserById mengembalikan pointer ke struct UserCore yang berisi data user
// berdasarkan id yang dikirimkan sebagai parameter.
// Fungsi ini mengembalikan error jika terjadi kesalahan saat query ke database.
func (u *userService) SelectUserById(id string) (*users.UserCore, error) {
	// Membuat query ke database untuk mengambil data user berdasarkan id.
	user, err := u.userData.SelectUserById(id)
	if err != nil {
		// Jika terjadi error saat query maka kembalikan error
		// dengan menggabungkan pesan error yang diterima.
		return nil, fmt.Errorf("gagal mengambil data user: %w", err)
	}

	// Kembalikan data user yang diperoleh dari repository.
	return user, nil
}

// UpdateUser implements users.ServiceUserInterface.
// Fungsi ini digunakan untuk mengupdate data user berdasarkan id yang diberikan.
// Fungsi ini menerima parameter input yang berisi data user yang ingin diupdate.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan saat query ke database.
func (u *userService) UpdateUser(input *users.UserCore, id string) error {
	if u == nil {
		// Jika objek userService adalah nil maka kembalikan error.
		return errors.New("user service: Nil service")
	}
	if u.userData == nil {
		// Jika repository userData adalah nil maka kembalikan error.
		return errors.New("user service: Nil Repository")
	}
	// Ambil data lama dari database berdasarkan id
	exisData, err := u.SelectUserById(id)
	if err != nil {
		// Jika terjadi error saat mengambil data maka kembalikan error.
		// Jika data tidak ditemukan, maka kembalikan error.
		if err == pgx.ErrNoRows {
			return errors.New("guru service: Data tidak ditemukan")
		}
		return fmt.Errorf("Id salah atau gagal mengambil data lama: %w", err)
	}

	// Gabungkan data baru dengan data lama
	// Jika field baru kosong, gunakan field dari data lama
	if input.Username == "" {
		// Jika field username kosong maka gunakan field username dari data lama.
		input.Username = exisData.Username
	}

	if input.Email == "" {
		// Jika field email kosong maka gunakan field email dari data lama.
		input.Email = exisData.Email
	}

	if input.Password == "" {
		// Jika field password kosong maka gunakan field password dari data lama.
		input.Password = exisData.Password
	}

	if input.Role == "" {
		// Jika field role kosong maka gunakan field role dari data lama.
		input.Role = exisData.Role
	}
	// Lakukan update data ke database
	if err := u.userData.UpdateUser(input, id); err != nil {
		// Jika terjadi error saat update maka kembalikan error.
		return err
	}

	// Kembalikan nil jika berhasil mengupdate data.
	return nil
}

// DeleteUserById mengimplementasikan users.ServiceUserInterface.
// Fungsi ini digunakan untuk menghapus data guru berdasarkan ID yang dikirimkan.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan saat menghapus data.
func (u *userService) DeleteUserById(id string) error {
	if u == nil {
		// Jika service kosong maka kembalikan error.
		return errors.New("user service: Nil service")
	}
	if u.userData == nil {
		// Jika repository kosong maka kembalikan error.
		return errors.New("user service: Nil Repository")
	}
	// Panggil fungsi DeleteUserById pada repository untuk menghapus data guru.
	// Jika terjadi error maka kembalikan error.
	if err := u.userData.DeleteUserById(id); err != nil {
		return fmt.Errorf("gagal menghapus data guru: %w", err)
	}
	return nil
}
