package users

import "time"

type (
	// UserCore merepresentasikan data user di database.
	// Data user terdiri atas:
	// 1. ID (string) sebagai id data user
	// 2. Username (string) sebagai nama pengguna
	// 3. Email (string) sebagai alamat email user
	// 4. Password (string) sebagai password user
	// 5. Role (string) sebagai peran user
	// 6. Update_At (time.Time) sebagai waktu update data user
	// 7. Delete_At (*time.Time) sebagai waktu delete data user
	UserCore struct {
		ID        string     `json:"id"`        // ID data user
		Username  string     `json:"username"`  // Nama pengguna
		Email     string     `json:"email"`     // Email user
		Password  string     `json:"password"`  // Password user
		Role      string     `json:"role"`      // Peran user
		Update_At time.Time  `json:"update_at"` // Waktu update data user
		Delete_At *time.Time `json:"delete_at"` // Waktu delete data user
	}

	// DataUserInterface merepresentasikan interface untuk data auth.
	// Interface ini digunakan untuk menghandle data auth yang berhubungan dengan user.
	DataUserInterface interface {
		// SelectAllUser mengembalikan slice users.UserCore yang berisi semua data user
		// yang tersimpan di database.
		// Fungsi ini mengembalikan error jika terjadi kesalahan saat query ke database.
		SelectAllUser() ([]UserCore, error)

		// SelectUserById mengembalikan pointer ke struct UserCore yang berisi data user
		// berdasarkan id yang dikirimkan sebagai parameter.
		// Fungsi ini mengembalikan error jika terjadi kesalahan saat query ke database.
		SelectUserById(id string) (*UserCore, error)

		// InsertUser implements users.DataUserInterface.
		// Fungsi ini digunakan untuk menginsert data user ke dalam database.
		// Fungsi ini menerima parameter input yang berisi data user yang ingin diinsert.
		// Fungsi ini akan mengembalikan error jika terjadi kesalahan saat query ke database.
		InsertUser(input *UserCore) error

		// UpdateUser implements users.DataUserInterface.
		// Fungsi ini digunakan untuk mengupdate data user berdasarkan ID yang diberikan.
		// Fungsi ini menerima parameter input yang berisi data user yang ingin diupdate.
		// Fungsi ini akan mengembalikan error jika terjadi kesalahan saat query ke database.
		UpdateUser(input *UserCore, id string) error

		// DeleteUserById mengimplementasikan users.DataUserInterface.
		// Fungsi ini digunakan untuk menghapus data guru berdasarkan ID yang dikirimkan.
		// Fungsi ini akan mengembalikan error jika terjadi kesalahan saat menghapus data.
		DeleteUserById(id string) error
	}

	// ServiceUserInterface merepresentasikan interface untuk service user.
	// Interface ini digunakan untuk menghandle data auth yang berhubungan dengan user.
	ServiceUserInterface interface {
		// SelectAllUser mengembalikan slice users.UserCore yang berisi semua data user
		// yang tersimpan di database.
		// Fungsi ini mengembalikan error jika terjadi kesalahan saat query ke database.
		SelectAllUser() ([]UserCore, error)

		// SelectUserById mengembalikan pointer ke struct UserCore yang berisi data user
		// berdasarkan id yang dikirimkan sebagai parameter.
		// Fungsi ini mengembalikan error jika terjadi kesalahan saat query ke database.
		SelectUserById(id string) (*UserCore, error)

		// InsertUser implements users.ServiceUserInterface.
		// Fungsi ini digunakan untuk menginsert data user ke dalam database.
		// Fungsi ini menerima parameter input yang berisi data user yang ingin diinsert.
		// Fungsi ini akan mengembalikan error jika terjadi kesalahan saat query ke database.
		InsertUser(input *UserCore) error

		// UpdateUser implements users.ServiceUserInterface.
		// Fungsi ini digunakan untuk mengupdate data user berdasarkan ID yang diberikan.
		// Fungsi ini menerima parameter input yang berisi data user yang ingin diupdate.
		// Fungsi ini akan mengembalikan error jika terjadi kesalahan saat query ke database.
		UpdateUser(input *UserCore, id string) error

		// DeleteUserById mengimplementasikan users.ServiceUserInterface.
		// Fungsi ini digunakan untuk menghapus data guru berdasarkan ID yang dikirimkan.
		// Fungsi ini akan mengembalikan error jika terjadi kesalahan saat menghapus data.
		DeleteUserById(id string) error
	}
)
