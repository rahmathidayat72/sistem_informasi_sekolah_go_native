package helper

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

type MetaToken struct {
	ID  string `json:"id"`
	Exp string `json:"exp"`
}

type AccessToken struct {
	Claims MetaToken
}

// SignToken digunakan untuk membuat token JWT baru berdasarkan data yang diberikan.
// Fungsi ini mengembalikan token yang ditandatangani, waktu kedaluwarsa, dan error jika ada.
func SignToken(data map[string]interface{}) (string, time.Time, error) {
	// Menetapkan waktu kedaluwarsa token secara hardcode
	expiryTime := time.Now().UTC().Add(time.Hour * 24) // Waktu kedaluwarsa 24 jam di UTC

	// Membuat klaim untuk token
	claims := jwt.MapClaims{}
	claims["exp"] = expiryTime.Unix() // Menetapkan waktu kedaluwarsa ke klaim

	// Menambahkan data tambahan ke klaim
	for key, value := range data {
		claims[key] = value
	}

	// Membuat token baru dengan metode penandatanganan HS256 dan klaim yang telah dibuat
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Menandatangani token dengan secret key dari environment variable JWT_SECRET
	accessToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", time.Time{}, err
	}

	// Mengembalikan token yang ditandatangani, waktu kedaluwarsa, dan nil untuk error
	return accessToken, expiryTime, nil
}

// Middleware untuk memverifikasi token JWT
// Middleware ini akan memverifikasi apakah token yang dikirimkan lewat header Authorization
// valid dan sesuai dengan secret key yang diatur di environment variable JWT_SECRET
// Jika token tidak valid maka akan dikembalikan error 401 Unauthorized
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Mendapatkan header Authorization dari request
		authHeader := r.Header.Get("Authorization")

		// Mendapatkan token dari header Authorization
		tokenString := GetTokenFromAuthorizationHeader(authHeader)

		// Jika token tidak ditemukan maka akan dikembalikan error 401 Unauthorized
		if tokenString == "" {
			JSONResponse(w, http.StatusUnauthorized, "Authorization token required")
			return
		}

		// Memverifikasi token yang diterima
		// Jika token tidak valid maka akan dikembalikan error 401 Unauthorized
		_, err := VerifyTokenHeader(tokenString)
		if err != nil {
			JSONResponse(w, http.StatusUnauthorized, "Invalid token: "+err.Error())
			return
		}

		// Jika token valid maka akan dilanjutkan ke handler berikutnya
		next.ServeHTTP(w, r)
	}
}

// Verifikasi token JWT yang diterima dari header Authorization
// Token yang diterima harus sesuai dengan secret key yang diatur di environment variable JWT_SECRET
// Jika token tidak valid maka akan dikembalikan error
func VerifyTokenHeader(requestToken string) (MetaToken, error) {
	// Buat token JWT baru
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		// Gunakan secret key yang diatur di environment variable JWT_SECRET
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	// Jika terjadi error maka akan dikembalikan error
	if err != nil {
		log.Println(err)
		return MetaToken{}, err
	}

	// Jika token tidak valid maka akan dikembalikan error
	if !token.Valid {
		log.Println("Token tidak valid")
		return MetaToken{}, jwt.ErrSignatureInvalid
	}

	// Dekode token menjadi MetaToken
	claimToken := DecodeToken(token)
	return claimToken.Claims, nil
}

func VerifyToken(accessToken string) (*jwt.Token, error) {
	jwtSecretKey := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	if !token.Valid {
		logrus.Error("Token is not valid")
		return nil, jwt.ErrSignatureInvalid
	}

	return token, nil
}

func DecodeToken(accessToken *jwt.Token) AccessToken {
	var token AccessToken
	stringify, err := json.Marshal(&accessToken)
	if err != nil {
		return token
	}
	err = json.Unmarshal(stringify, &token)
	if err != nil {
		return token
	}
	return token
}

// Fungsi untuk mengambil token dari header Authorization
func GetTokenFromAuthorizationHeader(authorizationHeader string) string {
	parts := strings.Split(authorizationHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}
