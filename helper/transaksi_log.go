package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TransactionLog merepresentasikan data log transaksi yang disimpan di database.
type TransactionLog struct {
	ID           int       `json:"id"`
	Timestamp    time.Time `json:"timestamp"`
	UserID       string    `json:"user_id"`
	Perangkat    string    `json:"perangkat"`
	ServiceName  string    `json:"service_name"`
	RequestBody  string    `json:"request_body"`
	ResponseBody string    `json:"response_body"`
	RequestParam string    `json:"request_param"`
	Result       string    `json:"result"`
	Header       string    `json:"header"`
}

// responseCategory adalah ResponseWriter custom untuk menangkap response body & status
type responseCategory struct {
	http.ResponseWriter
	body       *[]byte
	statusCode int
	wrote      bool // untuk mencegah WriteHeader dipanggil 2x
}

func (rw *responseCategory) WriteHeader(status int) {
	if rw.wrote {
		return
	}
	rw.statusCode = status
	rw.ResponseWriter.WriteHeader(status)
	rw.wrote = true
}

func (rw *responseCategory) Write(b []byte) (int, error) {
	*rw.body = append(*rw.body, b...)
	// default status jika belum pernah ditulis
	if !rw.wrote {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

// GetServiceNameFromEndpoint digunakan untuk memberi nama service berdasarkan endpoint
func GetServiceNameFromEndpoint(endpoint string) string {
	endpoint = strings.Split(endpoint, "?")[0]
	parts := strings.Split(strings.Trim(endpoint, "/"), "/")

	if len(parts) >= 2 {
		return parts[len(parts)-2] + "/" + parts[len(parts)-1]
	} else if len(parts) == 1 {
		return parts[0]
	}
	return ""
}

// LoggingMiddleware mengumpulkan data transaksi setiap request
func LoggingMiddleware(next http.Handler, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestBody []byte
		var responseBody []byte

		if db == nil {
			http.Error(w, "DB pool belum diinisialisasi", http.StatusInternalServerError)
			return
		}

		// baca request body
		if r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				requestBody = bodyBytes
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// Ambil userID dari token / body
		authHeader := r.Header.Get("Authorization")
		accessToken := GetTokenFromAuthorizationHeader(authHeader)

		var userID string
		if accessToken != "" {
			token, err := VerifyToken(accessToken)
			if err != nil {
				http.Error(w, "Token tidak valid", http.StatusUnauthorized)
				return
			}
			log.Println("Token:", token)

			metaToken, err := VerifyTokenHeader(accessToken)
			if err != nil {
				http.Error(w, "Header token tidak valid", http.StatusUnauthorized)
				return
			}
			userID = metaToken.ID
		} else {
			var bodyMap map[string]interface{}
			if err := json.Unmarshal(requestBody, &bodyMap); err == nil {
				if email, ok := bodyMap["email"].(string); ok {
					uid, err := GetUserIDByEmail(db, email)
					if err == nil {
						userID = uid
					}
				}
			}
		}

		// ResponseWriter custom
		responseWriter := &responseCategory{
			ResponseWriter: w,
			body:           &responseBody,
			statusCode:     200,
		}

		// Jalankan handler berikutnya
		next.ServeHTTP(responseWriter, r)

		// Simpan log di background
		go func() {
			perangkat, _ := os.Hostname()

			paramJSON, _ := json.Marshal(r.URL.Query())
			headerJSON, _ := json.Marshal(r.Header)

			// masking request body sensitif
			requestStr := "{}"
			if len(requestBody) > 0 {
				requestStr = maskSensitiveData(string(requestBody))
			}

			// response body → pastikan JSON valid
			responseStr := "{}"
			if json.Valid(responseBody) {
				responseStr = string(responseBody)
			} else if len(responseBody) > 0 {
				tmp, _ := json.Marshal(string(responseBody))
				responseStr = string(tmp)
			}

			resultStatus := "Success"
			if responseWriter.statusCode >= 400 {
				resultStatus = "Failed"
			}

			_, err := db.Exec(
				context.Background(),
				`INSERT INTO transaction_logs 
					(timestamp, user_id, perangkat, service_name, request_body, response_body, request_param, result, header)
				VALUES ($1, $2, $3, $4, $5::jsonb, $6::jsonb, $7::jsonb, $8, $9::jsonb)`,
				time.Now(),
				userID,
				perangkat,
				GetServiceNameFromEndpoint(r.RequestURI),
				requestStr,
				responseStr,
				string(paramJSON),
				resultStatus,
				string(headerJSON),
			)
			if err != nil {
				log.Println("[ERROR] Gagal simpan log:", err)
			}
		}()
	})
}

// GetUserIDByEmail ambil ID user berdasarkan email
func GetUserIDByEmail(db *pgxpool.Pool, email string) (string, error) {
	var userID string
	query := "SELECT id FROM users WHERE email = $1"
	err := db.QueryRow(context.Background(), query, email).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

// maskSensitiveData → sembunyikan field sensitif
func maskSensitiveData(data string) string {
	var bodyMap map[string]interface{}
	if err := json.Unmarshal([]byte(data), &bodyMap); err != nil {
		return data
	}

	sensitiveFields := []string{"password", "access_token", "refresh_token"}
	for _, field := range sensitiveFields {
		if val, ok := bodyMap[field].(string); ok {
			if field == "password" {
				bodyMap[field] = HashPassword(val)
			} else {
				bodyMap[field] = "***MASKED***"
			}
		}
	}

	hashedData, _ := json.Marshal(bodyMap)
	return string(hashedData)
}
