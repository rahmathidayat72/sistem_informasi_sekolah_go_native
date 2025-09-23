package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/kelas"
	"go_rest_native_sekolah/helper"
	"log"
	"net/http"
	"strings"
)

// KelasController merepresentasikan controller untuk menghandle HTTP request
// terkait data kelas. Struct ini mengandung field yang digunakan untuk
// mengakses service kelas.
type KelasController struct {
	// KelasService adalah field yang digunakan untuk mengakses service kelas.
	// Service ini berfungsi sebagai lapisan logika bisnis yang memungkinkan
	// controller untuk berinteraksi dengan data kelas di database.
	KelasService kelas.ServiceKelasInterface
}

// NewKelasController membuat objek KelasController dengan parameter KelasService.
// KelasController digunakan untuk menghandle logika bisnis yang berhubungan dengan data kelas.
// Jika parameter KelasService nil maka akan terjadi panic.
func NewKelasController(service kelas.ServiceKelasInterface) *KelasController {
	// Membuat objek KelasController baru dan menginisialisasi field KelasService
	return &KelasController{
		KelasService: service, // Menyimpan service kelas ke dalam field KelasService
	}
}

// Insert digunakan untuk menghandle HTTP request untuk menginsert data kelas ke dalam database.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (kc *KelasController) Insert(w http.ResponseWriter, r *http.Request) error {
	// Cek apakah kelas controller dan service kelas tidak nil.
	if kc == nil || kc.KelasService == nil {
		return errors.New("Nil controller")
	}

	// Deklarasikan objek yang digunakan untuk mengubah data inputan menjadi objek kelas core.
	var kelasReq KelasFormatter
	var kelasCore kelas.KelasCore

	// Dapatkan tipe konten yang dikirimkan oleh client.
	contentType := r.Header.Get("Content-Type")

	// Jika tipe konten adalah JSON, maka decode JSON ke dalam objek kelasReq.
	if strings.HasPrefix(contentType, "application/json") {
		err := json.NewDecoder(r.Body).Decode(&kelasReq)
		if err != nil {
			http.Error(w, "gagal membaca JSON", http.StatusBadRequest)
			return fmt.Errorf("error decoding JSON: %v", err)
		}
	} else {
		// Jika tipe konten bukan JSON maka decode data form ke dalam objek kelasReq.
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "gagal membaca form data", http.StatusBadRequest)
			return fmt.Errorf("error parsing form-data: %v", err)
		}
		kelasReq = KelasFormatter{
			Kelas:     r.FormValue("kelas"),
			ID_Guru:   r.FormValue("id_guru"),
			Nama_Guru: r.FormValue("nama_guru"),
		}
	}

	// Ubah objek kelasReq ke dalam objek kelasCore.
	kelasCore = FormatKelasRequestToCore(kelasReq)

	// Insert data kelas ke dalam database menggunakan service kelas.
	err := kc.KelasService.Insert(&kelasCore)
	if err != nil {
		http.Error(w, "gagal menyimpan data kelas", http.StatusInternalServerError)
		return fmt.Errorf("kelas controller: gagal insert kelas: %v", err)
	}

	// Pastikan respons mencerminkan data terbaru.
	formatKelas := FormatKelasList([]kelas.KelasCore{kelasCore})

	// Kirim response JSON.
	respon := helper.APIResponse(http.StatusCreated, "success insert kelas", formatKelas)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(respon)
	if err != nil {
		return fmt.Errorf("kelas controller: error encoding response: %v", err)
	}
	return nil
}

// Kelas digunakan untuk menghandle HTTP request untuk mengambil semua data kelas.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (kc *KelasController) Kelas(w http.ResponseWriter, r *http.Request) error {
	// Mengambil semua data kelas dari database melalui service kelas.
	kelas, err := kc.KelasService.SelectAll()
	if err != nil {
		// Jika terjadi error maka akan mengembalikan error dengan pesan "Error retrieving data".
		return fmt.Errorf("kelas controller: Error retrieving data: %v", err)
	}

	// Membuat slice KelasFormatter dan mengisi dengan data kelas yang diambil dari database.
	// KelasFormatter digunakan untuk memformat data kelas agar sesuai dengan kebutuhan response API.
	formatedKelas := FormatKelasList(kelas)

	// Membuat response API berupa JSON dengan status OK dan pesan "Success".
	// KelasFormatter yang diambil dari database akan di-encode menjadi JSON dan dikirim sebagai response.
	response := helper.APIResponse(http.StatusOK, "Berhasil mengambil data kelas dari database", formatedKelas)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// Jika terjadi error saat encoding response maka akan mengembalikan error dengan pesan "Error encoding response".
		return fmt.Errorf("kelas controller: Error encoding response: %v", err)
	}
	return nil // Mengembalikan nil karena tidak ada error
}

// GetKelasById digunakan untuk menghandle HTTP request untuk mengambil data kelas berdasarkan ID.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (kc *KelasController) GetKelasById(w http.ResponseWriter, r *http.Request) error {
	// Ambil ID kelas dari parameter query (?id=)
	id := r.URL.Query().Get("id")
	if id == "" {
		// Jika ID tidak ditemukan, kembalikan error dengan status 400.
		http.Error(w, "parameter 'id' wajib diisi", http.StatusBadRequest)
		return fmt.Errorf("kelas controller: ID kelas tidak ditemukan dalam query parameter")
	}

	// Panggil service untuk mengambil data kelas berdasarkan ID
	kelasData, err := kc.KelasService.SelectById(id)
	if err != nil {
		// Jika terjadi error saat mengambil data kelas, kembalikan error dengan pesan yang sesuai.
		return fmt.Errorf("kelas controller: gagal mengambil data kelas berdasarkan ID: %v", err)
	}

	// Format data menjadi list agar bisa diproses FormatKelasList
	formatedKelas := FormatKelasList([]kelas.KelasCore{*kelasData})

	// Buat response API
	// response ini berisi data kelas yang telah di-format dan di-encode menjadi JSON.
	// Jika terjadi error saat encoding maka kembalikan error dengan status 500.
	response := helper.APIResponse(http.StatusOK, "Success", formatedKelas)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// Jika terjadi error saat encoding response maka kembalikan error dengan pesan "Error encoding response".
		return fmt.Errorf("kelas controller: Error encoding response: %v", err)
	}

	// Jika tidak ada error maka kembalikan nil.
	return nil
}

// UpdateKelas digunakan untuk menghandle HTTP request untuk memperbarui data kelas berdasarkan ID yang dikirimkan.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (kc *KelasController) UpdateKelas(w http.ResponseWriter, r *http.Request) error {
	// Ambil ID kelas dari parameter query (?id=)
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// Jika tidak ada parameter 'id' maka kembalikan error dengan status 400.
		http.Error(w, "parameter 'id' wajib diisi", http.StatusBadRequest)
		return fmt.Errorf("kelas controller: ID kelas tidak ditemukan dalam query parameter")
	}

	// Log permintaan update untuk ID tertentu.
	log.Printf("Request update kelas dengan ID: %s", idStr)

	// Dekode request body menjadi objek KelasFormatter.
	var kelasReq KelasFormatter
	err := json.NewDecoder(r.Body).Decode(&kelasReq)
	if err != nil {
		// Jika terjadi error saat decoding maka kembalikan error dengan status 400.
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Gagal memproses data input", http.StatusBadRequest)
		return err
	}

	// Format data KelasFormatter menjadi objek KelasCore.
	updateKelas := FormatKelasRequestToCore(kelasReq)

	// Panggil service untuk memperbarui data kelas berdasarkan ID.
	err = kc.KelasService.Update(&updateKelas, idStr)
	if err != nil {
		// Jika terjadi error saat memperbarui data kelas maka kembalikan error sesuai dengan status error.
		// Jika terjadi error validasi, kirimkan response dengan status 400.
		if strings.Contains(err.Error(), "validation") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return err
		}

		// Jika ID kelas tidak ditemukan atau tidak ada perubahan data, kirimkan response dengan status 404.
		if strings.Contains(err.Error(), "tidak ditemukan") {
			http.Error(w, "ID kelas tidak ditemukan atau tidak ada perubahan data", http.StatusNotFound)
			return err
		}

		// Jika terjadi kesalahan umum saat memperbarui data, kirimkan response dengan status 500.
		http.Error(w, "Terjadi kesalahan saat memperbarui data", http.StatusInternalServerError)
		return err
	}

	// Ambil data kelas yang telah diupdate dari database.
	kelasUpdate, err := kc.KelasService.SelectById(idStr)
	if err != nil {
		// Jika terjadi error saat mengambil data kelas maka kembalikan error dengan status 500.
		http.Error(w, "gagal mengambil data kelas", http.StatusInternalServerError)
		return fmt.Errorf("kelas controller: gagal mengambil data kelas: %v", err)
	}

	// Format data kelas yang telah diupdate menjadi slice KelasCore.
	formatKelas := FormatKelasList([]kelas.KelasCore{*kelasUpdate})

	// Buat response API
	// response ini berisi data kelas yang telah di-format dan di-encode menjadi JSON.
	// Jika terjadi error saat encoding maka kembalikan error dengan status 500.
	respon := helper.APIResponse(http.StatusOK, "success update kelas", formatKelas)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(respon)
	if err != nil {
		return fmt.Errorf("kelas controller: error encoding response: %v", err)
	}

	// Jika tidak ada error maka kembalikan nil.
	return nil
}

// DeleteKelas digunakan untuk menghandle HTTP request untuk menghapus data kelas berdasarkan ID yang dikirimkan.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (kc *KelasController) DeleteKelas(w http.ResponseWriter, r *http.Request) error {
	// Ambil ID kelas dari parameter query (?id=)
	id := r.URL.Query().Get("id")
	if id == "" {
		// Jika ID tidak ditemukan, kembalikan error dengan status 400.
		http.Error(w, "parameter 'id' wajib diisi", http.StatusBadRequest)
		return fmt.Errorf("kelas controller: ID kelas tidak ditemukan dalam query parameter")
	}

	// Panggil service untuk menghapus data kelas berdasarkan ID
	// Jika terjadi error saat menghapus data kelas, kembalikan error dengan pesan yang sesuai.
	err := kc.KelasService.DeleteById(id)
	if err != nil {
		http.Error(w, "gagal menghapus data kelas", http.StatusInternalServerError)
		return fmt.Errorf("kelas controller: gagal menghapus data kelas: %v", err)
	}

	// Buat response API
	// response ini berisi pesan "success deleted kelas id: " dan status OK.
	// Jika terjadi error saat encoding maka kembalikan error dengan status 500.
	respon := helper.APIResponse(http.StatusOK, "success deleted kelas id: "+id, nil)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(respon)
	if err != nil {
		return fmt.Errorf("kelas controller: error encoding response: %v", err)
	}

	// Jika tidak ada error maka kembalikan nil.
	return nil
}
