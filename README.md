# 📘 Sistem Informasi Sekolah (Go Native + PostgreSQL)

## 📝 Deskripsi

Aplikasi sistem informasi sekolah untuk mengelola:

- Data **User**
- Data **Siswa**
- Data **Guru**
- Data **Kelas**
- Data **Mata Pelajaran**

### Fitur utama

- CRUD (Create, Read, Update, Delete) siswa, guru, kelas, mapel dan user
- Autentikasi **JWT**
- Logging transaksi request/response

---

## ⚙️ Tech Stack

- **Bahasa**: ![Go](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white&style=for-the-badge) (Native, tanpa framework HTTP tambahan)
- **Database**: ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-336791?logo=postgresql&logoColor=white&style=for-the-badge)
- **Library**:
  - [`github.com/jackc/pgx/v5 v5.7.5`](https://pkg.go.dev/github.com/jackc/pgx/v5) → driver PostgreSQL
  - [`github.com/joho/godotenv v1.5.1`](https://pkg.go.dev/github.com/joho/godotenv) → load file `.env`
  <!-- - [`github.com/stretchr/testify v1.10.0`](https://pkg.go.dev/github.com/stretchr/testify) → testing -->

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
   APP_PORT=1234 (hanya contoh)

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

### 👤 User

- GET /users → Ambil semua user

- GET /users/userbyid?{id} → Ambil detail user

- POST /users/tambah → Tambah user

- PUT /users/update?{id} → Update user

- DELETE /users/deleted?{id} → Hapus user

### 👨‍🏫 Guru

- GET /guru → list semua guru

- POST /guru/tambah → tambah guru

- GET /guru/gurubyid?{id} → detail guru

- PUT /guru/update?{id} → update guru

- DELETE /guru/deleted?{id} → hapus guru

### 👨‍🎓 Siswa

- GET /siswa → list semua siswa

- POST /siswa/tambah → tambah siswa

- GET /siswa/siswabyid?{id} → detail siswa

- PUT /siswa/update?{id} → update siswa

- DELETE /siswa/deleted?{id} → hapus siswa

### 🏫 Kelas

- GET /kelas → list semua kelas

- POST /kelas/tambah → tambah kelas

- GET /kelas/kelasbyid?{id} → detail kelas

- PUT /kelas/update?{id} → update kelas

- DELETE /kelas/deleted?{id} → hapus kelas

### 📖 Mata Pelajaran

- GET /mapel → list semua mapel

- POST /mapel/tambah → tambah mapel

- GET /mapel/update?{id} → detail mapel

- PUT /mapel/update?{id} → update mapel

- DELETE /mapel/deleted?{id} → hapus mapel

---

## ✨ Catatan

- Pastikan file .env sesuai konfigurasi database lokal.

- Endpoint membutuhkan JWT Token setelah login.

- Logging transaksi otomatis tersimpan di tabel transaction_logs.

---

## 👨‍💻 Kontributor

- **Rahmat Hidayat** – Developer utama
