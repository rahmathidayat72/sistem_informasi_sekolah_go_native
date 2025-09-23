package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	matapelajaran "go_rest_native_sekolah/features/mata_pelajaran"
	"go_rest_native_sekolah/helper"
	"log"
	"net/http"
	"strings"
)

// MataPelajaranController adalah struktur yang bertanggung jawab untuk
// menangani permintaan HTTP yang berkaitan dengan mata pelajaran.
type MataPelajaranController struct {
	// MataPelajaranService adalah antarmuka layanan yang digunakan untuk
	// berinteraksi dengan data mata pelajaran.
	MataPelajaranService matapelajaran.ServiceMapelInterface
}

// NewMataPelajaranController adalah fungsi yang digunakan untuk membuat
// instance dari MataPelajaranController. Fungsi ini akan mengembalikan
// instance dari MataPelajaranController yang memiliki property
// MataPelajaranService yang memiliki nilai service yang diinputkan.
//
// Fungsi ini digunakan untuk membuat instance dari MataPelajaranController
// yang siap digunakan untuk menangani permintaan HTTP yang berkaitan dengan
// data mata pelajaran.
func NewMataPelajaranController(service matapelajaran.ServiceMapelInterface) MataPelajaranController {
	// Buat instance dari MataPelajaranController dan set MataPelajaranService
	// dengan nilai service yang diinputkan.
	return MataPelajaranController{MataPelajaranService: service}
}

// InsertMapel digunakan untuk menghandle permintaan HTTP POST untuk
// menginsert data mata pelajaran ke dalam database.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (mpc *MataPelajaranController) InsertMapel(w http.ResponseWriter, r *http.Request) error {
	// Cek apakah controller dan service mata pelajaran tidak nil.
	if mpc == nil || mpc.MataPelajaranService == nil {
		return errors.New("Nil controller")
	}

	// Deklarasikan objek yang digunakan untuk mengubah data inputan menjadi objek mata pelajaran core.
	var mapelReq FormatterMataPelajaran
	var mapelCore matapelajaran.MataPelajaranCore

	// Dapatkan tipe konten yang dikirimkan oleh client.
	contentType := r.Header.Get("Content-Type")

	// Jika tipe konten adalah JSON, maka decode JSON ke dalam objek mapelReq.
	if strings.HasPrefix(contentType, "application/json") {
		err := json.NewDecoder(r.Body).Decode(&mapelReq)
		if err != nil {
			// Jika terjadi error maka kirimkan response error 400 Bad Request.
			http.Error(w, "gagal membaca JSON", http.StatusBadRequest)
			return fmt.Errorf("error decoding JSON: %v", err)
		}
	} else {
		// Jika tipe konten bukan JSON maka decode data form ke dalam objek mapelReq.
		err := r.ParseForm()
		if err != nil {
			// Jika terjadi error maka kirimkan response error 400 Bad Request.
			http.Error(w, "gagal membaca form data", http.StatusBadRequest)
			return fmt.Errorf("error parsing form-data: %v", err)
		}
		mapelReq = FormatterMataPelajaran{
			Nama_Pelajaran: r.FormValue("nama_pelajaran"),
			ID_Guru:        r.FormValue("id_guru"),
			Guru:           r.FormValue("guru"),
			Kelas_ID:       r.FormValue("kelas_id"),
			Nama_Kelas:     r.FormValue("nama_kelas"),
			Deskripsi:      r.FormValue("deskripsi"),
		}
	}

	// Ubah request ke Core
	mapelCore = FormatterMapelRequestToCore(mapelReq)

	// Insert data mata pelajaran ke dalam database menggunakan service mata pelajaran.
	err := mpc.MataPelajaranService.InsertMapel(&mapelCore)
	if err != nil {
		return err
	}

	// Pastikan respons mencerminkan data terbaru.
	formatMapel := FormatterMapelList([]matapelajaran.MataPelajaranCore{mapelCore})

	// Kirim response JSON.
	respon := helper.APIResponse(http.StatusCreated, "Success insert data matapelajaran", formatMapel)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(respon)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}
	return nil
}

// Mapel adalah fungsi yang digunakan untuk mengembalikan data mata pelajaran
// yang tersedia di database.
//
// Fungsi ini akan mengembalikan data mata pelajaran yang diambil dari database
// dalam bentuk JSON. Data yang diambil berupa nama mata pelajaran, ID guru, nama
// guru, ID kelas, nama kelas, dan deskripsi.
//
// Jika terjadi error saat mengambil data maka akan dikembalikan dalam bentuk
// response JSON dengan kode status 500 Internal Server Error.
func (mpc *MataPelajaranController) Mapel(w http.ResponseWriter, r *http.Request) error {
	if mpc == nil || mpc.MataPelajaranService == nil {
		// Jika controller atau service nil maka akan dikembalikan error.
		return errors.New("Nil controller")
	}
	mapel, err := mpc.MataPelajaranService.SelectAllMapel()
	// Panggil fungsi SelectAllMapel pada service untuk mengambil data mata pelajaran.
	if err != nil {
		// Jika terjadi error maka akan dikembalikan dalam bentuk response JSON.
		return err
	}
	formatMapel := FormatterMapelList(mapel)
	// Ubah data yang diambil menjadi format JSON yang dibutuhkan.
	respon := helper.APIResponse(http.StatusOK, "Success get data mapel", formatMapel)
	// Buatkan response JSON yang dibutuhkan.
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(respon)
	if err != nil {
		// Jika terjadi error saat mengencode JSON maka akan dikembalikan dalam bentuk
		// response JSON dengan kode status 500 Internal Server Error.
		return fmt.Errorf("error encoding JSON: %v", err)
	}
	return nil // Jika tidak ada error maka kembalikan nil
}

// GetMapelById adalah fungsi yang digunakan untuk mengembalikan data mata pelajaran
// yang dicari berdasarkan ID.
//
// Fungsi ini akan mengembalikan data mata pelajaran yang diambil dari database
// dalam bentuk JSON. Data yang diambil berupa nama mata pelajaran, ID guru, nama
// guru, ID kelas, nama kelas, dan deskripsi.
//
// Jika terjadi error saat mengambil data maka akan dikembalikan dalam bentuk
// response JSON dengan kode status 500 Internal Server Error.
func (mpc *MataPelajaranController) GetMapelById(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	// Ambil parameter 'id' yang dikirimkan lewat URL.
	if id == "" {
		http.Error(w, "parameter 'id' wajib diisi", http.StatusBadRequest)
		// Jika parameter 'id' kosong maka akan dikembalikan error dengan kode status 400 Bad Request.
		return errors.New("missing 'id' query parameter")
	}
	mapelData, err := mpc.MataPelajaranService.SelectMapelById(id)
	// Panggil fungsi SelectMapelById pada service untuk mengambil data mata pelajaran yang dicari.
	if err != nil {
		// Jika terjadi error maka akan dikembalikan dalam bentuk response JSON.
		return err
	}
	formatMapel := FormatterMapelList([]matapelajaran.MataPelajaranCore{*mapelData})
	// Ubah data yang diambil menjadi format JSON yang dibutuhkan.
	respon := helper.APIResponse(http.StatusOK, "Success get data mapelById", formatMapel)
	// Buatkan response JSON yang dibutuhkan.
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(respon)
	if err != nil {
		// Jika terjadi error saat mengencode JSON maka akan dikembalikan dalam bentuk
		// response JSON dengan kode status 500 Internal Server Error.
		return fmt.Errorf("error encoding JSON: %v", err)
	}
	return nil // Jika tidak ada error maka kembalikan nil
}

// UpdateMapel digunakan untuk mengupdate data mata pelajaran berdasarkan ID.
// Fungsi ini menerima request berupa JSON yang berisi nama mata pelajaran, ID guru,
// nama guru, ID kelas, nama kelas, dan deskripsi.
//
// Fungsi ini akan mengembalikan data mata pelajaran yang diupdate dalam bentuk JSON.
// Jika terjadi error maka akan dikembalikan dalam bentuk response JSON dengan kode
// status 500 Internal Server Error.
func (mpc *MataPelajaranController) UpdateMapel(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	// Ambil parameter 'id' yang dikirimkan lewat URL.
	if id == "" {
		http.Error(w, "parameter 'id' wajib diisi", http.StatusBadRequest)
		// Jika parameter 'id' kosong maka akan dikembalikan error dengan kode status 400 Bad Request.
		return errors.New("missing 'id' query parameter")
	}
	var mapelReq FormatterMataPelajaran
	// Deklarasikan objek yang digunakan untuk mengubah data inputan menjadi objek mata pelajaran core.
	err := json.NewDecoder(r.Body).Decode(&mapelReq)
	// Dekode data yang dikirimkan lewat body menjadi objek mata pelajaran.
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "gagal memproses data input", http.StatusBadRequest)
		// Jika terjadi error saat decoding maka akan dikembalikan error dengan kode status 400 Bad Request.
		return err
	}
	mapelUpdate := FormatterMapelRequestToCore(mapelReq)
	// Ubah data yang diambil menjadi format objek mata pelajaran core.
	err = mpc.MataPelajaranService.UpdateMapel(&mapelUpdate, id)
	// Panggil fungsi UpdateMapel pada service untuk mengupdate data mata pelajaran yang dicari.
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		// Jika terjadi error maka akan dikembalikan dalam bentuk response JSON dengan kode status 400 Bad Request.
		return err
	}
	mapelData, err := mpc.MataPelajaranService.SelectMapelById(id)
	// Panggil fungsi SelectMapelById pada service untuk mengambil data mata pelajaran yang diupdate.
	if err != nil {
		http.Error(w, "Gagal mengambil data setelah update", http.StatusInternalServerError)
		// Jika terjadi error maka akan dikembalikan dalam bentuk response JSON dengan kode status 500 Internal Server Error.
		return err
	}
	formatMapel := FormatterMapelList([]matapelajaran.MataPelajaranCore{*mapelData})
	// Ubah data yang diambil menjadi format JSON yang dibutuhkan.
	respon := helper.APIResponse(http.StatusOK, "Berhasil mengupdate data mapel", formatMapel)
	// Buatkan response JSON yang dibutuhkan.
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(respon)
	if err != nil {
		// Jika terjadi error saat mengencode JSON maka akan dikembalikan dalam bentuk
		// response JSON dengan kode status 500 Internal Server Error.
		return fmt.Errorf("error encoding JSON: %v", err)
	}
	return nil // Jika tidak ada error maka kembalikan nil
}

// DeleteMapel digunakan untuk menghandle permintaan HTTP DELETE untuk
// menghapus data mata pelajaran yang dicari berdasarkan id.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (mpc *MataPelajaranController) DeleteMapel(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	// Ambil parameter 'id' yang dikirimkan lewat URL.
	if id == "" {
		http.Error(w, "parameter 'id' wajib diisi", http.StatusBadRequest)
		// Jika parameter 'id' kosong maka akan dikembalikan error dengan kode status 400 Bad Request.
		return errors.New("missing 'id' query parameter")
	}
	err := mpc.MataPelajaranService.DeleteMapel(id)
	// Panggil fungsi DeleteMapel pada service untuk menghapus data mata pelajaran yang dicari.
	if err != nil {
		return err
		// Jika terjadi error maka akan dikembalikan dalam bentuk response JSON dengan kode status 500 Internal Server Error.
	}
	respon := helper.APIResponse(http.StatusOK, "Berhasil menghapus data mapel", nil)
	// Buatkan response JSON yang dibutuhkan.
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(respon)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
		// Jika terjadi error saat mengencode JSON maka akan dikembalikan dalam bentuk
		// response JSON dengan kode status 500 Internal Server Error.
	}
	return nil // Jika tidak ada error maka kembalikan nil
}
