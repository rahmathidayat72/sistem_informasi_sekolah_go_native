package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/siswa"
	"go_rest_native_sekolah/helper"
	"log"
	"net/http"
	"strings"
)

// SiswaController digunakan untuk menghandle HTTP request yang berhubungan dengan data siswa.
//
// Struktur ini memiliki 1 field, yaitu SiswaService yang memiliki tipe siswa.ServiceSiswaInterface.
// SiswaService digunakan untuk mengakses data siswa di database.
type SiswaController struct {
	// SiswaService digunakan untuk mengakses data siswa di database.
	// Field ini memiliki tipe siswa.ServiceSiswaInterface yang berisi method-method
	// untuk mengakses data siswa, seperti Insert, SelectAll, SelectById, Update, Delete.
	SiswaService siswa.ServiceSiswaInterface
}

// NewSiswaController membuat objek SiswaController baru dengan parameter service.
// SiswaController digunakan untuk menghandle logika bisnis yang berhubungan dengan data siswa.
// Jika parameter service nil maka akan terjadi panic.
func NewSiswaController(service siswa.ServiceSiswaInterface) *SiswaController {
	// Membuat objek SiswaController baru dan menginisialisasi field SiswaService
	return &SiswaController{
		SiswaService: service, // Menyimpan service siswa ke dalam field SiswaService
	}
}

// InsertSiswa digunakan untuk menghandle HTTP request POST untuk menginsert data siswa ke dalam database.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (sc *SiswaController) InsertSiswa(w http.ResponseWriter, r *http.Request) error {
	// Cek apakah controller tidak nil dan service siswa tidak nil.
	if sc == nil || sc.SiswaService == nil {
		return errors.New("Nil controller")
	}

	// Deklarasikan objek yang digunakan untuk mengubah data inputan menjadi objek siswa core.
	var siswaReq SiswaFormatter
	var siswaCore siswa.SiswaCore

	// Dapatkan tipe konten yang dikirimkan oleh client.
	contentType := r.Header.Get("Content-Type")

	// Jika tipe konten adalah JSON, maka decode JSON ke dalam objek siswaReq.
	if strings.HasPrefix(contentType, "application/json") {
		// Jika terjadi error maka kirimkan response error 400 Bad Request.
		err := json.NewDecoder(r.Body).Decode(&siswaReq)
		if err != nil {
			http.Error(w, "gagal membaca JSON", http.StatusBadRequest)
			return fmt.Errorf("error decoding JSON: %v", err)
		}
	} else {
		// Jika tipe konten bukan JSON maka decode data form ke dalam objek siswaReq.
		err := r.ParseForm()
		if err != nil {
			// Jika terjadi error maka kirimkan response error 400 Bad Request.
			http.Error(w, "gagal membaca form data", http.StatusBadRequest)
			return fmt.Errorf("error parsing form-data: %v", err)
		}
		siswaReq = SiswaFormatter{
			Nama:       r.FormValue("nama"),
			Kelas_ID:   r.FormValue("kelas_id"),
			Nama_Kelas: r.FormValue("nama_kelas"),
			Email:      r.FormValue("email"),
			Alamat:     r.FormValue("alamat"),
		}
	}

	// Ubah request ke Core
	siswaCore = FormatSiswaRequestToCore(siswaReq)

	// Insert ke service
	err := sc.SiswaService.InsertSiswa(&siswaCore)
	if err != nil {
		return err
	}

	// Pastikan respons mencerminkan data terbaru
	formatKelas := FormatterKelasList([]siswa.SiswaCore{siswaCore})

	// Kirim response JSON
	respon := helper.APIResponse(http.StatusCreated, "Success insert data siswa", formatKelas)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(respon)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}
	return nil
}

// Siswa digunakan untuk menghandle HTTP request GET untuk mengambil semua data siswa.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (sc *SiswaController) Siswa(w http.ResponseWriter, r *http.Request) error {
	// Cek apakah controller tidak nil dan service siswa tidak nil.
	if sc == nil || sc.SiswaService == nil {
		return errors.New("Nil controller")
	}

	// Panggil service untuk mengambil semua data siswa.
	siswa, err := sc.SiswaService.SelectAllSiswa()
	if err != nil {
		// Jika terjadi error saat mengambil data siswa, maka kembalikan error.
		return err
	}

	// Format data menjadi list agar bisa diproses FormatterKelasList.
	formatKelas := FormatterKelasList(siswa)

	// Buat response API.
	// Response ini berisi data siswa yang telah di-format dan di-encode menjadi JSON.
	// Jika terjadi error saat encoding maka kembalikan error dengan status 500.
	respon := helper.APIResponse(http.StatusOK, "Success get data siswa", formatKelas)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(respon)
	if err != nil {
		// Jika terjadi error saat encoding response maka kembalikan error dengan pesan "Error encoding response".
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	// Jika tidak ada error maka kembalikan nil.
	return nil
}

// GetSiswaById digunakan untuk menghandle HTTP request GET untuk mengambil data siswa berdasarkan ID.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (sc *SiswaController) GetSiswaById(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	// Ambil parameter 'id' yang dikirimkan lewat URL.
	if id == "" {
		// Jika parameter 'id' kosong maka akan dikembalikan error dengan kode status 400 Bad Request.
		http.Error(w, "parameter 'id' wajib diisi", http.StatusBadRequest)
		return errors.New("missing 'id' query parameter")
	}
	siswaData, err := sc.SiswaService.SelectById(id)
	// Panggil service untuk mengambil data siswa berdasarkan ID.
	if err != nil {
		// Jika terjadi error saat mengambil data siswa, maka kembalikan error.
		return err
	}
	formatKelas := FormatterKelasList([]siswa.SiswaCore{*siswaData})
	// Format data menjadi list agar bisa diproses FormatterKelasList.
	respon := helper.APIResponse(http.StatusOK, "Success get data siswaById", formatKelas)
	// Buat response API.
	// Response ini berisi data siswa yang telah di-format dan di-encode menjadi JSON.
	// Jika terjadi error saat encoding maka kembalikan error dengan status 500.
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(respon)
	if err != nil {
		// Jika terjadi error saat encoding response maka kembalikan error dengan pesan "Error encoding response".
		return fmt.Errorf("error encoding JSON: %v", err)
	}
	return nil // Jika tidak ada error maka kembalikan nil
}

// UpdateSiswa digunakan untuk menghandle HTTP request untuk memperbarui data siswa berdasarkan ID yang dikirimkan.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (sc *SiswaController) UpdateSiswa(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	// Ambil parameter 'id' yang dikirimkan lewat URL.
	// Jika parameter 'id' kosong maka akan dikembalikan error dengan kode status 400 Bad Request.
	if id == "" {
		http.Error(w, "parameter 'id' wajib diisi", http.StatusBadRequest)
		return errors.New("missing 'id' query parameter")
	}

	var siswaReq SiswaFormatter
	// Dekode request body menjadi objek SiswaFormatter.
	// Jika terjadi error maka kembalikan error dengan status 400 Bad Request.
	err := json.NewDecoder(r.Body).Decode(&siswaReq)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "gagal memproses data input", http.StatusBadRequest)
		return err
	}

	siswaUpdate := FormatSiswaRequestToCore(siswaReq)
	// Format data SiswaFormatter menjadi objek SiswaCore.
	// Jika terjadi error maka kembalikan error.
	err = sc.SiswaService.Update(&siswaUpdate, id)
	if err != nil {
		// Jika terjadi error saat memperbarui data siswa, maka kembalikan error.
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	// Ambil data siswa yang telah diupdate dari database
	siswaData, err := sc.SiswaService.SelectById(id)
	if err != nil {
		// Jika terjadi error saat mengambil data siswa maka kembalikan error dengan status 500
		http.Error(w, "Gagal mengambil data setelah update", http.StatusInternalServerError)
		return err
	}
	// Format data siswa yang diupdate menjadi objek FormatterKelasList
	formatKelas := FormatterKelasList([]siswa.SiswaCore{*siswaData})
	// Buat response API
	// response ini berisi pesan "Berhasil mengupdate data siswa" dan data siswa yang diupdate
	// Jika terjadi error saat encoding maka kembalikan error dengan status 500
	respon := helper.APIResponse(http.StatusOK, "Berhasil mengupdate data siswa", formatKelas)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(respon)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}
	// Jika tidak ada error maka kembalikan nil
	return nil
}

// DeleteSiswa digunakan untuk menghandle HTTP request DELETE untuk menghapus data siswa berdasarkan ID yang dikirimkan.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (sc *SiswaController) DeleteSiswa(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	// Ambil parameter 'id' yang dikirimkan lewat URL.
	// Jika parameter 'id' kosong maka akan dikembalikan error dengan kode status 400 Bad Request.
	if id == "" {
		http.Error(w, "parameter 'id' wajib diisi", http.StatusBadRequest)
		return errors.New("missing 'id' query parameter")
	}
	err := sc.SiswaService.DeleteById(id)
	// Panggil service untuk menghapus data siswa berdasarkan ID.
	// Jika terjadi error saat menghapus data siswa, maka kembalikan error.
	if err != nil {
		return err
	}
	respon := helper.APIResponse(http.StatusOK, "Berhasil menghapus data siswa", nil)
	// Buat response API.
	// Response ini berisi pesan "Berhasil menghapus data siswa" dan status OK.
	// Jika terjadi error saat encoding maka kembalikan error dengan status 500.
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(respon)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}
	return nil
}
