package helper

import "golang.org/x/crypto/bcrypt"

// HashPassword digunakan untuk mengenkripsi password menjadi bentuk hash.
// Fungsi ini membutuhkan parameter berupa string yang berisi password yang akan dienkripsi.
// Fungsi ini mengembalikan string yang berisi hash dari password.
func HashPassword(pass string) string {
	password := []byte(pass)
	// GenerateFromPassword digunakan untuk mengenerate hash dari password yang diinput.
	// Fungsi ini membutuhkan parameter berupa slice of bytes yang berisi password dan cost yang berisi tingkat kesulitan enkripsi.
	// Fungsi ini mengembalikan slice of bytes yang berisi hash dari password dan error jika terjadi kesalahan.
	hassPwd, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		// Jika terjadi kesalahan maka panic dengan error yang terjadi.
		panic(err)
	}
	// Fungsi ini mengembalikan string yang berisi hash dari password.
	return string(hassPwd)
}

// CheckPassword digunakan untuk membandingkan password yang diinput dengan hash yang diinput.
// Fungsi ini membutuhkan parameter berupa string yang berisi password yang diinput dan string yang berisi hash yang diinput.
// Fungsi ini mengembalikan nilai boolean yang berisi hasil perbandingan.
func CheckPassword(password, hash string) bool {
	// passwordBytes digunakan untuk mengkonversi string password menjadi slice of bytes.
	passwordBytes := []byte(password)

	// hashBytes digunakan untuk mengkonversi string hash menjadi slice of bytes.
	hashBytes := []byte(hash)

	// CompareHashAndPassword digunakan untuk membandingkan password yang diinput dengan hash yang diinput.
	// Fungsi ini membutuhkan parameter berupa slice of bytes yang berisi hash dan slice of bytes yang berisi password yang diinput.
	// Fungsi ini mengembalikan error jika terjadi kesalahan.
	err := bcrypt.CompareHashAndPassword(hashBytes, passwordBytes)

	// Jika error-nya nil maka artinya password yang diinput sama dengan hash yang diinput.
	// Fungsi ini mengembalikan nilai boolean yang berisi hasil perbandingan.
	return err == nil
}
