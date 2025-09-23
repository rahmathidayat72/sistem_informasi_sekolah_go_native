package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/guru"
	"go_rest_native_sekolah/helper"
	"log"
	"net/http"
	"strings"
)

// Gurucontroller merepresentasikan controller untuk tabel guru.
// guruService digunakan untuk menghandle logika bisnis yang berhubungan dengan tabel guru.
type Gurucontroller struct {
	guruService guru.ServiceGuruInterface
}

// NewGuruController membuat objek guruController dengan parameter guruService.
// guruController digunakan untuk menghandle logika bisnis yang berhubungan dengan tabel guru.
// Jika parameter guruService nil maka akan terjadi panic.
func NewGuruController(service guru.ServiceGuruInterface) *Gurucontroller {
	// if service == nil {
	// 	panic("guru controller: Nil service")
	// }
	return &Gurucontroller{guruService: service}
}

// Guru digunakan untuk menghandle HTTP request untuk mengambil semua data guru.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (gc *Gurucontroller) Guru(w http.ResponseWriter, r *http.Request) error {
	// Mengambil semua data guru dari database melalui service guru.
	gurus, err := gc.guruService.GetAllGuru()
	if err != nil {
		// Jika terjadi error maka akan mengembalikan error dengan pesan "Error retrieving data".
		return fmt.Errorf("guru controller: Error retrieving data: %v", err)
	}

	// Membuat slice GuruFormatter dan mengisi dengan data guru yang diambil dari database.
	// GuruFormatter digunakan untuk memformat data guru agar sesuai dengan kebutuhan response API.
	formattedGurus := FormatGuruList(gurus) // format sebelum dikirim

	// Membuat response API berupa JSON dengan status OK dan pesan "Success".
	// GuruFormatter yang diambil dari database akan di-encode menjadi JSON dan dikirim sebagai response.
	response := helper.APIResponse(http.StatusOK, "Berhasil mengambil data guru dari database", formattedGurus)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// Jika terjadi error saat encoding response maka akan mengembalikan error dengan pesan "Error encoding response".
		return fmt.Errorf("guru controller: Error encoding response: %v", err)
	}

	return nil
}

// InsertGuru digunakan untuk menghandle HTTP request untuk menginsert data guru ke dalam database.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (gc *Gurucontroller) InsertGuru(w http.ResponseWriter, r *http.Request) error {
	// Cek service tidak nil.
	if gc == nil || gc.guruService == nil {
		return fmt.Errorf("guru controller: service is nil")
	}

	var guruReq GuruFormatter  // objek GuruFormatter digunakan untuk mengubah data guru yang diinputkan menjadi objek GuruCore.
	var guruCore guru.GuruCore // objek GuruCore digunakan untuk mengakses data guru yang diinputkan.

	contentType := r.Header.Get("Content-Type") // Ambil tipe konten yang dikirimkan oleh client.

	if strings.HasPrefix(contentType, "application/json") {
		// JSON body
		err := json.NewDecoder(r.Body).Decode(&guruReq)
		if err != nil {
			http.Error(w, "gagal membaca JSON", http.StatusBadRequest)
			return fmt.Errorf("error decoding JSON: %v", err)
		}
	} else {
		// Form-data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "gagal membaca form data", http.StatusBadRequest)
			return fmt.Errorf("error parsing form-data: %v", err)
		}
		guruReq = GuruFormatter{
			ID_User: r.FormValue("id_user"),
			Nama:    r.FormValue("nama"),
			Email:   r.FormValue("email"),
			Alamat:  r.FormValue("alamat"),
		}
	}

	// Validasi awal.
	if guruReq.Nama == "" || guruReq.Email == "" || guruReq.Alamat == "" {
		http.Error(w, "nama, email, dan alamat wajib diisi", http.StatusBadRequest)
		return fmt.Errorf("validasi gagal: nama, email, dan alamat harus diisi")
	}

	// Format ke core.
	guruCore = FormatGuruRequestToCore(guruReq)

	// Simpan data.
	err := gc.guruService.InsertGuru(&guruCore)
	if err != nil {
		http.Error(w, "gagal menyimpan data guru", http.StatusInternalServerError)
		return fmt.Errorf("gagal insert guru: %v", err)
	}

	// Format response.
	formattedGurus := FormatGuruList([]guru.GuruCore{guruCore})
	response := helper.APIResponse(http.StatusOK, "Berhasil menginsert data guru ke database", formattedGurus)

	// Tulis response.
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return fmt.Errorf("error encoding response: %v", err)
	}

	return nil
}

// UpdateGuru digunakan untuk menghandle HTTP request untuk memperbarui data guru berdasarkan ID yang dikirimkan.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (gc *Gurucontroller) UpdateGuru(w http.ResponseWriter, r *http.Request) error {
	// Ambil ID guru dari parameter query (?id=)
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// Jika tidak ada parameter 'id' maka kembalikan error dengan status 400.
		http.Error(w, "parameter 'id' wajib diisi", http.StatusBadRequest)
		return errors.New("missing 'id' query parameter")
	}
	log.Printf("Request update guru dengan ID: %s", idStr)

	// Dekode request body menjadi objek GuruFormatter.
	var guruReq GuruFormatter
	err := json.NewDecoder(r.Body).Decode(&guruReq)
	if err != nil {
		// Jika terjadi error saat decoding maka kembalikan error dengan status 400.
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Gagal memproses data input", http.StatusBadRequest)
		return err
	}

	// Format data GuruFormatter menjadi objek GuruCore.
	updateGuru := FormatGuruRequestToCore(guruReq)

	// Panggil service untuk memperbarui data guru berdasarkan ID.
	err = gc.guruService.UpdateGuru(&updateGuru, idStr)
	if err != nil {
		// Jika terjadi error saat memperbarui data guru maka kembalikan error sesuai dengan status error.
		if strings.Contains(err.Error(), "validation") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return err
		}
		if strings.Contains(err.Error(), "tidak ditemukan") {
			http.Error(w, "ID guru tidak ditemukan atau tidak ada perubahan data", http.StatusNotFound)
			return err
		}
		http.Error(w, "Terjadi kesalahan saat memperbarui data", http.StatusInternalServerError)
		return err
	}

	// Ambil data guru yang telah diupdate dari database.
	updatedGuru, err := gc.guruService.SelectById(idStr)
	if err != nil {
		// Jika terjadi error saat mengambil data guru maka kembalikan error dengan status 500.
		http.Error(w, "Gagal mengambil data setelah update", http.StatusInternalServerError)
		return err
	}

	// Format data guru menjadi slice GuruCore.
	formattedGurus := FormatGuruList([]guru.GuruCore{*updatedGuru})

	// Buat response API dengan status OK dan pesan "Berhasil mengupdate data guru ke database".
	response := helper.APIResponse(http.StatusOK, "Berhasil mengupdate data guru ke database", formattedGurus)

	// Set header Content-Type menjadi application/json.
	w.Header().Set("Content-Type", "application/json")

	// Encode response API menjadi JSON dan tulis ke response writer.
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Jika terjadi error saat encoding maka kembalikan error dengan status 500.
		return fmt.Errorf("error encoding response: %v", err)
	}
	// Jika tidak ada error maka kembalikan nil.
	return nil
}

// GetGuruById digunakan untuk menghandle HTTP request untuk mengambil data guru berdasarkan ID.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (gc *Gurucontroller) GetGuruById(w http.ResponseWriter, r *http.Request) error {
	// Ambil ID guru dari parameter query (?id=)
	id := r.URL.Query().Get("id")
	if id == "" {
		// Jika ID tidak ditemukan, kembalikan error dengan status 400.
		http.Error(w, "parameter 'id' wajib diisi", http.StatusBadRequest)
		return fmt.Errorf("guru controller: ID guru tidak ditemukan dalam query parameter")
	}

	// Panggil service untuk mengambil data guru berdasarkan ID
	guruData, err := gc.guruService.SelectById(id)
	if err != nil {
		// Jika terjadi error saat mengambil data guru, kembalikan error dengan pesan yang sesuai.
		return fmt.Errorf("guru controller: gagal mengambil data guru berdasarkan ID: %v", err)
	}

	// Format data menjadi list agar bisa diproses FormatGuruList
	formattedGurus := FormatGuruList([]guru.GuruCore{*guruData})

	// Buat response API
	// response ini berisi data guru yang telah di-format dan di-encode menjadi JSON.
	// Jika terjadi error saat encoding maka kembalikan error dengan status 500.
	response := helper.APIResponse(http.StatusOK, "Berhasil mengambil data guru berdasarkan ID", formattedGurus)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return fmt.Errorf("guru controller: error saat encoding response: %v", err)
	}

	// Jika tidak ada error maka kembalikan nil.
	return nil
}

// DeleteGuru digunakan untuk menghandle HTTP request untuk menghapus data guru berdasarkan ID yang dikirimkan.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (gc *Gurucontroller) DeleteGuru(w http.ResponseWriter, r *http.Request) error {
	// Ambil ID guru dari parameter query (?id=)
	// Jika ID tidak ditemukan, kembalikan error dengan status 400.
	id := r.URL.Query().Get("id")
	if id == "" {
		return fmt.Errorf("guru controller: ID guru tidak ditemukan dalam query parameter")
	}

	// Panggil service untuk menghapus data guru berdasarkan ID
	// Jika terjadi error saat menghapus data guru, kembalikan error dengan pesan yang sesuai.
	err := gc.guruService.DeleteById(id)
	if err != nil {
		return fmt.Errorf("guru controller: gagal menghapus data guru berdasarkan ID: %v", err)
	}

	// Buat response API
	// response ini berisi pesan "Terhapus" dan status OK.
	// Jika terjadi error saat encoding maka kembalikan error dengan status 500.
	response := helper.APIResponse(http.StatusOK, "Berhasil menghapus data guru berdasarkan ID", "Terhapus")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return fmt.Errorf("guru controller: error saat encoding response: %v", err)
	}

	// Jika tidak ada error maka kembalikan nil.
	return nil
}
