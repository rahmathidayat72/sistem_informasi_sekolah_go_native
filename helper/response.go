package helper

import (
	"encoding/json"
	"net/http"
)

// Response adalah struktur data yang digunakan untuk mengembalikan respon ke client.
// Struktur ini terdiri dari:
//   - Message: pesan yang diberikan dalam respon.
//   - Code: kode status HTTP yang diberikan dalam respon.
//   - Success: boolean yang menunjukan apakah respon berhasil atau tidak.
//   - Data: data opsional yang diberikan dalam respon.
type Response struct {
	Message string      `json:"message"` // Pesan yang diberikan dalam respon.
	Code    int         `json:"code"`    // Kode status HTTP yang diberikan dalam respon.
	Success bool        `json:"success"` // Boolean yang menunjukan apakah respon berhasil atau tidak.
	Data    interface{} `json:"data,omitempty"`
}

// APIResponse membuat struktur respon untuk permintaan API.
// Fungsi ini mengambil kode status HTTP, pesan, dan data opsional,
// dan mengembalikan struktur Response.
//
// Parameter:
//
//	status: kode status HTTP yang akan digunakan dalam respon.
//	message: pesan yang akan digunakan dalam respon.
//	data: data opsional yang akan digunakan dalam respon.
//
// Nilai kembali:
//
//	struktur Response yang berisi informasi yang diberikan.
func APIResponse(status int, message string, data interface{}) Response {
	// Tentukan kesuksesan berdasarkan kode status
	success := status < 400

	// Jika data nil, jangan set apapun (agar omitempty berfungsi)
	var response Response
	if data == nil {
		response = Response{
			Message: message,
			Code:    status,
			Success: success,
			// Data tidak diset = tetap nil = tidak muncul di JSON
		}
	} else {
		response = Response{
			Message: message,
			Code:    status,
			Success: success,
			Data:    data,
		}
	}

	return response
}

// JSONResponse mengirimkan respon dalam format JSON.
// Fungsi ini mengambil responsewriter, kode status HTTP, dan data yang akan dikirimkan,
// dan mengirimkan respon dalam format JSON.
//
// Parameter:
//
//	w: responsewriter yang akan digunakan untuk mengirimkan respon.
//	status: kode status HTTP yang akan digunakan dalam respon.
//	payload: data yang akan dikirimkan dalam respon.
func JSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	// Set Content-Type header ke application/json
	w.Header().Set("Content-Type", "application/json")

	// Set kode status HTTP yang akan digunakan dalam respon
	w.WriteHeader(status)

	// Encode data yang akan dikirimkan dalam format JSON
	json.NewEncoder(w).Encode(payload)
}
