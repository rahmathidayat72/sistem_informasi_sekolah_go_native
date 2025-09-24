package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DBPool *pgxpool.Pool

// InitPostgreSQLPool digunakan untuk menginisialisasi pool koneksi ke database PostgreSQL.
// Fungsi ini akan mengembalikan pointer ke objek pgxpool.Pool dan error.
// Jika terjadi error maka akan mengembalikan error dengan pesan "Gagal koneksi ke database".
// Fungsi ini juga akan mengisi variabel DBPool dengan pointer ke objek pgxpool.Pool yang diinisialisasi.
// DBPool dapat digunakan untuk mengakses database tanpa harus membuat objek pgxpool.Pool lagi.
func InitPostgreSQLPool() (*pgxpool.Pool, error) {
	// // Memuat file .env
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatalf("Gagal memuat file .env: %v", err)
	// }
	
	// Load env sesuai APP_ENV (.env.development / .env.production / .env.testing)
	LoadEnv()

	// Mengambil environment variable
	host := os.Getenv("DBHOST")
	port := os.Getenv("DBPORT")
	user := os.Getenv("DBUSER")
	password := os.Getenv("DBPASS")
	dbname := os.Getenv("DBNAME")

	// Membuat string DSN untuk koneksi ke database PostgreSQL
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&timezone=Asia/Shanghai", user, password, host, port, dbname)

	// Konfigurasi pool
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Gagal parsing konfigurasi pool: %v", err)
	}

	// Konfigurasi pool
	config.MaxConns = 500                      // Jumlah koneksi maksimum
	config.MinConns = 10                       // Jumlah koneksi minimum
	config.MaxConnLifetime = 5 * time.Second   // Masa hidup koneksi maksimum
	config.HealthCheckPeriod = 2 * time.Second // Periode pengecekan kesehatan

	// Membuat context dengan timeout 5 detik
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Membuat pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("gagal koneksi ke database: %w", err)
	}

	// Mengecek koneksi ke database
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping ke database gagal: %w", err)
	}
	// Menyimpan pointer ke objek pgxpool.Pool yang diinisialisasi ke dalam variabel DBPool
	// Variabel DBPool dapat digunakan untuk mengakses database tanpa harus membuat objek pgxpool.Pool lagi
	DBPool = pool

	// Mencetak pesan informasi bahwa berhasil koneksi ke database PostgreSQL
	log.Println("[INFO] Berhasil koneksi ke database PostgreSQL")

	// Mengembalikan pointer ke objek pgxpool.Pool yang diinisialisasi dan nilai error yang nil
	return pool, nil
}
