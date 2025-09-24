# 📘 Sistem Informasi Sekolah (Go Native + PostgreSQL)

## 📝 Deskripsi

Aplikasi sistem informasi sekolah untuk mengelola:

- Data **Siswa**
- Data **Guru**
- Data **Kelas**
- Data **Mata Pelajaran**

### Fitur utama

- CRUD (Create, Read, Update, Delete) siswa, guru, kelas, mapel
- Autentikasi **JWT**
- Logging transaksi request/response

---

## ⚙️ Tech Stack

- **Bahasa**: Go (Native, tanpa framework HTTP tambahan)
- **Database**: PostgreSQL
- **Library**:
  - [`github.com/jackc/pgx/v5 v5.7.5`](https://pkg.go.dev/github.com/jackc/pgx/v5) → driver PostgreSQL
  - [`github.com/joho/godotenv v1.5.1`](https://pkg.go.dev/github.com/joho/godotenv) → load file `.env`
  - [`github.com/stretchr/testify v1.10.0`](https://pkg.go.dev/github.com/stretchr/testify) → testing

---

## 📂 Struktur Proyek

```
├── config/ # Konfigurasi .env dan database
├── features/ # Modul (guru, siswa, kelas, mapel, auth, users)
│ ├── controllers/ # Controller tiap modul
│ ├── model/ # Model & query database
│ ├── service/ # Business logic
│ └── entities.go # Struct entity
├── helper/ # Helper umum (hash, middleware, logging)
├── router/ # Routing endpoint
├── db.txt # File DDL untuk membuat database
├── .exp.env # Contoh konfigurasi environment
├── main.go # Entry point aplikasi
└── README.md # Dokumentasi proyek
```

---

## 🚀 Instalasi & Menjalankan

1. Clone repository:
   ```bash
   git clone https://github.com/rahmathidayat72/sistem_informasi_sekolah_go_native.git
   cd sistem_informasi_sekolah_go_native
   ```
2. Buat file .env dengan menyalin isi dari .exp.env, lalu isi sesuai konfigurasi komputer masing-masing.
   Contoh .exp.env:

   ```
   APP_NAME="Sistem Informasi Sekolah"
   APP_PORT=8081

   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=yourpassword
   DB_NAME=sekolah_db

   ```

3. Buat database sesuai schema di file db.txt
4. Jalankan aplikasi:
   ```
   go run main.go
   ```
5. Server akan berjalan sesuai port pada .env, misalnya:
   ```
   🌐 Server berjalan di http://localhost:your_number_port
   ```

---

## 📡 API Endpoint

### 🔑 Auth

- POST /login → login & dapatkan JWT

### 👨‍🏫 Guru

- GET /guru → list semua guru

- POST /guru → tambah guru

- GET /guru/{id} → detail guru

- PUT /guru/{id} → update guru

- DELETE /guru/{id} → hapus guru

### 👨‍🎓 Siswa

- GET /siswa → list semua siswa

- POST /siswa → tambah siswa

- GET /siswa/{id} → detail siswa

- PUT /siswa/{id} → update siswa

- DELETE /siswa/{id} → hapus siswa

### 🏫 Kelas

- GET /kelas → list semua kelas

- POST /kelas → tambah kelas

- GET /kelas/{id} → detail kelas

- PUT /kelas/{id} → update kelas

- DELETE /kelas/{id} → hapus kelas

### 📖 Mata Pelajaran

- GET /mapel → list semua mapel

- POST /mapel → tambah mapel

- GET /mapel/{id} → detail mapel

- PUT /mapel/{id} → update mapel

- DELETE /mapel/{id} → hapus mapel

---

## ✨ Catatan

- Pastikan file .env sesuai konfigurasi database lokal.

- Endpoint membutuhkan JWT Token setelah login.

- Logging transaksi otomatis tersimpan di tabel transaction_logs.

---

## 👨‍💻 Kontributor

- Rahmat Hidayat – Developer utama
