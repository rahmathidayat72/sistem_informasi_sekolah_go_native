package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/users"
	"go_rest_native_sekolah/helper"
	"log"
	"net/http"
	"strings"
)

// UserController adalah struktur data yang berisi property userService.
// Property ini memiliki nilai berupa interface ServiceUserInterface yang
// digunakan untuk mengakses data user.
type UserController struct {
	// userService adalah property yang berisi interface untuk mengakses data user.
	userService users.ServiceUserInterface
}

// NewUsesController adalah fungsi yang digunakan untuk membuat instance baru dari UserController.
// Fungsi ini menerima parameter userService yang merupakan interface untuk mengakses data user.
func NewUsesController(userService users.ServiceUserInterface) *UserController {
	// Kembalikan pointer ke UserController baru dengan menginisialisasi property userService.
	return &UserController{
		userService: userService, // Inisialisasi property userService dengan nilai dari parameter.
	}
}

// Users adalah fungsi yang digunakan untuk mengambil data user.
// Fungsi ini akan mengembalikan data user dalam bentuk JSON.
func (uc *UserController) Users(w http.ResponseWriter, r *http.Request) error {
	// Jika controller adalah nil, maka kita akan mengembalikan error.
	if uc == nil {
		http.Error(w, "Service is not available", http.StatusInternalServerError)
		return fmt.Errorf("user controller: controller is nil")
	}

	// Jika service user adalah nil, maka kita akan mengembalikan error.
	if uc.userService == nil {
		http.Error(w, "Service is not available", http.StatusInternalServerError)
		return fmt.Errorf("user controller: service is nil")
	}

	// Ambil data user dari database.
	users, err := uc.userService.SelectAllUser()
	if err != nil {
		// Jika ada error, maka kita akan mengembalikan error.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("user controller: error retrieving users: %v", err)
	}

	// Jika data user tidak ditemukan, maka kita akan mengembalikan error.
	if len(users) == 0 {
		http.Error(w, "Data not found", http.StatusNotFound)
		return fmt.Errorf("user controller: data not found")
	}

	// Format data user menjadi bentuk JSON.
	formattedUsers := FormatUserList(users)

	// Buatkan response JSON yang berisi data user.
	response := helper.APIResponse(http.StatusOK, "Success", formattedUsers)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Jika ada error dalam mengencode JSON, maka kita akan mengembalikan error.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("user controller: error encoding response: %v", err)
	}

	// Kembalikan nilai error yang kosong.
	return nil
}

// InsertUser adalah fungsi yang digunakan untuk menghandle HTTP request
// untuk menginsert data user ke dalam database.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (uc *UserController) InsertUser(w http.ResponseWriter, r *http.Request) error {
	// Jika controller atau service user adalah nil, maka kita akan mengembalikan error.
	if uc == nil || uc.userService == nil {
		return fmt.Errorf("user controller: service is nil")
	}

	// Deklarasikan objek yang digunakan untuk mengubah data inputan menjadi objek user core.
	var userReq UserFormatter
	var usersCore users.UserCore

	// Dapatkan tipe konten yang dikirimkan oleh client.
	contentType := r.Header.Get("Content-Type")

	// Jika tipe konten adalah JSON, maka decode JSON ke dalam objek userReq.
	if strings.HasPrefix(contentType, "application/json") {
		// Jika JSON body
		err := json.NewDecoder(r.Body).Decode(&userReq)
		if err != nil {
			// Jika terjadi error, maka kita akan mengembalikan error.
			http.Error(w, "gagal membaca JSON", http.StatusBadRequest)
			return fmt.Errorf("error decoding JSON: %v", err)
		}
	} else {
		// Jika form-data
		err := r.ParseForm()
		if err != nil {
			// Jika terjadi error, maka kita akan mengembalikan error.
			http.Error(w, "gagal membaca form data", http.StatusBadRequest)
			return fmt.Errorf("error parsing form-data: %v", err)
		}
		// Jika tipe konten bukan JSON, maka decode data form ke dalam objek userReq.
		userReq = UserFormatter{
			Username: r.FormValue("username"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
			Role:     r.FormValue("role"),
		}
	}

	// Validasi awal.
	// Jika field username, email, password, atau role kosong, maka kita akan mengembalikan error.
	if userReq.Username == "" || userReq.Email == "" || userReq.Password == "" || userReq.Role == "" {
		http.Error(w, "semua field wajib diisi", http.StatusBadRequest)
		return fmt.Errorf("validasi gagal: field kosong")
	}

	// Format ke core.
	// Ubah objek userReq ke dalam objek userCore.
	usersCore = FormatUserRequestToCore(userReq)

	// Simpan data user ke dalam database menggunakan service user.
	err := uc.userService.InsertUser(&usersCore)
	if err != nil {
		// Jika terjadi error, maka kita akan mengembalikan error.
		http.Error(w, "gagal menyimpan data user", http.StatusInternalServerError)
		return fmt.Errorf("user controller: gagal insert user: %v", err)
	}

	// Format respons tanpa merubah nilai usersCore.
	// Buatkan response JSON yang berisi data user yang telah diinsert.
	formattedUser := FormatUserList([]users.UserCore{usersCore})
	response := helper.APIResponse(http.StatusCreated, "User created", formattedUser)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// Jika terjadi error dalam mengencode JSON, maka kita akan mengembalikan error.
		return fmt.Errorf("user controller: error encoding response: %v", err)
	}

	// Kembalikan nilai error yang kosong.
	return nil
}

// GetUserById digunakan untuk menghandle HTTP request untuk mengambil data user berdasarkan ID.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (uc *UserController) GetUserById(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	// Ambil parameter 'id' yang dikirimkan lewat URL.
	if id == "" {
		// Jika parameter 'id' kosong maka akan dikembalikan error dengan kode status 400 Bad Request.
		http.Error(w, "parameter 'id' wajib diisi", http.StatusBadRequest)
		return errors.New("missing 'id' query parameter")
	}
	userData, err := uc.userService.SelectUserById(id)
	// Panggil service untuk mengambil data user berdasarkan ID.
	if err != nil {
		// Jika terjadi error saat mengambil data user, maka kembalikan error.
		return fmt.Errorf("user controller: gagal mengambil data user berdasarkan ID: %v", err)
	}
	formatuser := FormatUserList([]users.UserCore{*userData})
	// Format data menjadi list agar bisa diproses FormatUserList.
	response := helper.APIResponse(http.StatusOK, "Success", formatuser)
	// Buat response API.
	// Response ini berisi data user yang telah di-format dan di-encode menjadi JSON.
	// Jika terjadi error saat encoding maka kembalikan error dengan status 500.
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// Jika terjadi error saat encoding response maka kembalikan error dengan pesan "Error encoding response".
		return fmt.Errorf("user controller: error encoding response: %v", err)
	}
	return nil // Jika tidak ada error maka kembalikan nil.
}

// UpdateUser digunakan untuk menghandle HTTP request untuk memperbarui data user berdasarkan ID yang dikirimkan.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// Jika parameter 'id' kosong maka akan dikembalikan error dengan kode status 400 Bad Request.
		http.Error(w, "parameter 'id' wajib diisi", http.StatusBadRequest)
		return errors.New("missing 'id' query parameter")
	}
	// Log permintaan update untuk ID tertentu.
	log.Printf("Request update user dengan ID: %s", idStr)

	// Deklarasikan objek yang digunakan untuk mengubah data inputan menjadi objek user core.
	var userReq UserFormatter

	// Dekode request body menjadi objek userReq.
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		// Jika terjadi error saat decoding maka kembalikan error dengan status 400.
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Gagal memproses data input", http.StatusBadRequest)
		return err
	}

	// Format data userReq menjadi objek user core.
	updateUser := FormatUserRequestToCore(userReq)

	// Panggil service untuk memperbarui data user berdasarkan ID.
	err = uc.userService.UpdateUser(&updateUser, idStr)
	if err != nil {
		// Jika terjadi error saat memperbarui data user maka kembalikan error sesuai dengan status error.
		if strings.Contains(err.Error(), "validation") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return err
		}
		if strings.Contains(err.Error(), "tidak ditemukan") {
			http.Error(w, "ID user tidak ditemukan atau tidak ada perubahan data", http.StatusNotFound)
			return err
		}
		http.Error(w, "Terjadi kesalahan saat memperbarui data", http.StatusInternalServerError)
		return err
	}

	// Ambil data user yang telah diupdate dari database.
	updatedUser, err := uc.userService.SelectUserById(idStr)
	if err != nil {
		// Jika terjadi error saat mengambil data user maka kembalikan error dengan status 500.
		http.Error(w, "Gagal mengambil data setelah update", http.StatusInternalServerError)
		return err
	}

	// Format data menjadi list agar bisa diproses FormatUserList.
	formattedUser := FormatUserList([]users.UserCore{*updatedUser})

	// Buat response API.
	// Response ini berisi data user yang telah di-format dan di-encode menjadi JSON.
	// Jika terjadi error saat encoding maka kembalikan error dengan status 500.
	response := helper.APIResponse(http.StatusOK, "Success", formattedUser)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		// Jika terjadi error saat encoding response maka kembalikan error dengan pesan "Error encoding response".
		return fmt.Errorf("user controller: error encoding response: %v", err)
	}
	return nil // Jika tidak ada error maka kembalikan nil.
}

// DeleteUser digunakan untuk menghandle HTTP request DELETE untuk menghapus data user berdasarkan ID yang dikirimkan.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	// Ambil parameter 'id' yang dikirimkan lewat URL.
	// Jika parameter 'id' kosong maka akan dikembalikan error dengan kode status 400 Bad Request.
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "parameter 'id' wajib diisi", http.StatusBadRequest)
		return errors.New("missing 'id' query parameter")
	}
	// Panggil service untuk menghapus data user berdasarkan ID.
	// Jika terjadi error saat menghapus data user maka kembalikan error.
	err := uc.userService.DeleteUserById(id)
	if err != nil {
		return fmt.Errorf("user controller: gagal menghapus data user berdasarkan ID: %v", err)
	}
	// Buat response API.
	// Response ini berisi pesan "success deleted user id: "+id dan status OK.
	// Jika terjadi error saat encoding maka kembalikan error dengan status 500.
	response := helper.APIResponse(http.StatusOK, "success deleted user id: "+id, nil)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return fmt.Errorf("user controller: error encoding response: %v", err)
	}
	// Jika tidak ada error maka kembalikan nil.
	return nil
}
