package model

import (
	"context"
	"errors"
	"fmt"
	matapelajaran "go_rest_native_sekolah/features/mata_pelajaran"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// mataPelajaranQuery adalah struktur data yang berisi pointer ke objek pgxpool.Pool.
// Struktur data ini digunakan untuk menghandle query ke database yang berhubungan dengan tabel mata_pelajaran.
type mataPelajaranQuery struct {
	db *pgxpool.Pool // db adalah pointer ke objek pgxpool.Pool yang berisi koneksi database.
}

// NewDataMataPelajaran membuat objek mataPelajaranQuery yang berisi koneksi database.
// Fungsi ini digunakan untuk menginisialisasi objek mataPelajaranQuery yang berisi koneksi database.
// Jika parameter db nil maka akan terjadi panic.
// Fungsi ini digunakan untuk menghandle query ke database yang berhubungan dengan tabel mata_pelajaran.
func NewDataMataPelajaran(db *pgxpool.Pool) matapelajaran.DataMataPelajaranInterface {
	// Jika db nil maka akan terjadi panic
	if db == nil {
		panic("Nil database")
	}

	// Membuat objek mataPelajaranQuery baru dan menginisialisasi field db dengan parameter db
	return &mataPelajaranQuery{db: db}
}

// InsertMapel digunakan untuk menginsert data mata pelajaran ke dalam database.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (m *mataPelajaranQuery) InsertMapel(insert *matapelajaran.MataPelajaranCore) error {
	// Cek koneksi database, jika nil kembalikan error.
	if m.db == nil {
		return errors.New("Nil database")
	}

	// Cek apakah data yang akan diinsert nil, jika iya kembalikan error.
	if insert == nil {
		return errors.New("insert data is nil")
	}

	// Jika ID belum diisi, generate UUID baru.
	if insert.ID == "" {
		insert.ID = uuid.New().String()
	}

	// Validasi dan sinkronisasi Nama_Guru & ID_Guru.
	switch {
	case insert.Guru != "" && insert.ID_Guru == "":
		// Jika hanya Nama_Guru diisi, cari ID_Guru berdasarkan Nama_Guru.
		insert.Guru = strings.TrimSpace(insert.Guru)
		var guruID string
		err := m.db.QueryRow(context.Background(),
			"SELECT id FROM guru WHERE TRIM(nama) ILIKE TRIM($1)", insert.Guru).Scan(&guruID)
		if err != nil {
			log.Printf("InsertMapel: nama guru '%s' tidak ditemukan", insert.Guru)
			return fmt.Errorf("guru dengan nama '%s' tidak ditemukan", insert.Guru)
		}
		insert.ID_Guru = guruID

	case insert.ID_Guru != "" && insert.Guru == "":
		// Jika hanya ID_Guru diisi, cari Nama_Guru berdasarkan ID_Guru.
		var namaGuru string
		err := m.db.QueryRow(context.Background(),
			"SELECT nama FROM guru WHERE id = $1", insert.ID_Guru).Scan(&namaGuru)
		if err != nil {
			log.Printf("InsertMapel: ID guru '%s' tidak ditemukan", insert.ID_Guru)
			return fmt.Errorf("guru dengan ID '%s' tidak ditemukan", insert.ID_Guru)
		}
		insert.Guru = namaGuru

	case insert.ID_Guru != "" && insert.Guru != "":
		// Jika keduanya diisi, validasi apakah cocok.
		insert.Guru = strings.TrimSpace(insert.Guru)
		var existingName string
		err := m.db.QueryRow(context.Background(),
			"SELECT nama FROM guru WHERE id = $1", insert.ID_Guru).Scan(&existingName)
		if err != nil {
			log.Printf("InsertMapel: ID guru '%s' tidak ditemukan", insert.ID_Guru)
			return fmt.Errorf("guru dengan ID '%s' tidak ditemukan", insert.ID_Guru)
		}
		if strings.ToLower(strings.TrimSpace(existingName)) != strings.ToLower(insert.Guru) {
			log.Printf("InsertMapel: Nama guru tidak cocok. Dapat: '%s', seharusnya: '%s'", insert.Guru, existingName)
			return fmt.Errorf("nama guru '%s' tidak cocok dengan ID guru '%s'", insert.Guru, insert.ID_Guru)
		}
	}

	// Validasi dan sinkronisasi Nama_Kelas & Kelas_ID.
	switch {
	case insert.Nama_Kelas != "" && insert.Kelas_ID == "":
		// Jika hanya Nama_Kelas diisi, cari Kelas_ID berdasarkan Nama_Kelas.
		insert.Nama_Kelas = strings.TrimSpace(insert.Nama_Kelas)
		var kelasID string
		err := m.db.QueryRow(context.Background(),
			"SELECT id FROM kelas WHERE TRIM(kelas) ILIKE TRIM($1)", insert.Nama_Kelas).Scan(&kelasID)
		if err != nil {
			log.Printf("InsertMapel: nama kelas '%s' tidak ditemukan", insert.Nama_Kelas)
			return fmt.Errorf("kelas dengan nama '%s' tidak ditemukan", insert.Nama_Kelas)
		}
		insert.Kelas_ID = kelasID

	case insert.Kelas_ID != "" && insert.Nama_Kelas == "":
		// Jika hanya Kelas_ID diisi, cari Nama_Kelas berdasarkan Kelas_ID.
		var namaKelas string
		err := m.db.QueryRow(context.Background(),
			"SELECT kelas FROM kelas WHERE id = $1", insert.Kelas_ID).Scan(&namaKelas)
		if err != nil {
			log.Printf("InsertMapel: ID kelas '%s' tidak ditemukan", insert.Kelas_ID)
			return fmt.Errorf("kelas dengan ID '%s' tidak ditemukan", insert.Kelas_ID)
		}
		insert.Nama_Kelas = namaKelas

	case insert.Kelas_ID != "" && insert.Nama_Kelas != "":
		// Jika keduanya diisi, validasi apakah cocok.
		insert.Nama_Kelas = strings.TrimSpace(insert.Nama_Kelas)
		var existingKelas string
		err := m.db.QueryRow(context.Background(),
			"SELECT kelas FROM kelas WHERE id = $1", insert.Kelas_ID).Scan(&existingKelas)
		if err != nil {
			log.Printf("InsertMapel: ID kelas '%s' tidak ditemukan", insert.Kelas_ID)
			return fmt.Errorf("kelas dengan ID '%s' tidak ditemukan", insert.Kelas_ID)
		}
		if strings.ToLower(strings.TrimSpace(existingKelas)) != strings.ToLower(insert.Nama_Kelas) {
			log.Printf("InsertMapel: Nama kelas tidak cocok. Dapat: '%s', seharusnya: '%s'", insert.Nama_Kelas, existingKelas)
			return fmt.Errorf("nama kelas '%s' tidak cocok dengan ID kelas '%s'", insert.Nama_Kelas, insert.Kelas_ID)
		}
	}

	// Siapkan ID_Guru & Kelas_ID agar bisa null jika kosong.
	var idGuruParam interface{}
	if insert.ID_Guru == "" {
		idGuruParam = nil
	} else {
		idGuruParam = insert.ID_Guru
	}

	var idKelasParam interface{}
	if insert.Kelas_ID == "" {
		idKelasParam = nil
	} else {
		idKelasParam = insert.Kelas_ID
	}

	// Eksekusi query insert data mata pelajaran ke dalam database.
	_, err := m.db.Exec(context.Background(),
		"INSERT INTO mata_pelajaran (id, nama_pelajaran, id_guru, kelas_id, deskripsi) VALUES ($1, $2, $3, $4, $5)",
		insert.ID, insert.Nama_Pelajaran, idGuruParam, idKelasParam, insert.Deskripsi)
	if err != nil {
		log.Printf("InsertMapel error exec: %v", err)
		return fmt.Errorf("insert failed: %w", err)
	}

	// Kembalikan nil jika tidak terjadi kesalahan.
	return nil
}

// SelectAllMapel implements matapelajaran.DataMataPelajaranInterface.
// Fungsi ini digunakan untuk mengambil semua data mata pelajaran dari database.
// Fungsi ini mengembalikan slice MataPelajaranCore yang berisi data mata pelajaran.
// Jika terjadi error maka fungsi ini akan mengembalikan error.
func (m *mataPelajaranQuery) SelectAllMapel() ([]matapelajaran.MataPelajaranCore, error) {
	if m.db == nil {
		// Jika database tidak ada, kembalikan error.
		return nil, errors.New("Nil database")
	}
	// Buat query untuk mengambil semua data mata pelajaran.
	query := `SELECT 
    mp.id,
    mp.nama_pelajaran,
    mp.id_guru,
    g.nama AS nama_guru,
    mp.kelas_id,
    k.kelas AS nama_kelas,
    mp.deskripsi
FROM mata_pelajaran mp
LEFT JOIN guru g ON mp.id_guru = g.id
LEFT JOIN kelas k ON mp.kelas_id = k.id
WHERE mp.delete_at IS NULL;`

	// Jalankan query dan simpan hasilnya dalam rows.
	rows, err := m.db.Query(context.Background(), query)
	if err != nil {
		// Jika terjadi error saat eksekusi query, log error dan kembalikan.
		log.Printf("SelectAllMapel error exec: %v", err)
		return nil, fmt.Errorf("select failed: %w", err)
	}
	defer rows.Close() // Pastikan rows ditutup setelah selesai digunakan.

	// Deklarasi variabel untuk menyimpan hasil.
	var result []matapelajaran.MataPelajaranCore

	// Iterasi melalui hasil rows.
	// Fungsi ini digunakan untuk mengiterasi data yang diambil dari database dan
	// memasukkan data tersebut dalam slice MataPelajaranCore.
	for rows.Next() {
		var mp MataPelajaran // Deklarasi variabel untuk menyimpan data yang diiterasi.

		// Pindai setiap baris ke dalam variabel mp.
		// Fungsi Scan digunakan untuk memindai setiap baris yang diiterasi
		// dan menyimpannya dalam variabel mp.
		err = rows.Scan(&mp.ID, &mp.Nama_Pelajaran, &mp.ID_Guru, &mp.Guru, &mp.Kelas_ID, &mp.Nama_Kelas, &mp.Deskripsi)
		if err != nil {
			// Jika terjadi error saat scan, log error dan kembalikan.
			log.Printf("SelectAllMapel error scan: %v", err)
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		// Ubah data mp menjadi MataPelajaranCore dan tambahkan ke result.
		core := FormatterResponse(mp)
		result = append(result, core)
	}
	log.Printf("Successfully fetched %d mata pelajaran from database", len(result))
	// Kembalikan slice MataPelajaranCore yang berisi data mata pelajaran.
	return result, nil
}

// SelectMapelById implements matapelajaran.DataMataPelajaranInterface.
// Fungsi ini digunakan untuk mengambil data mata pelajaran berdasarkan id.
// Fungsi ini mengembalikan objek MataPelajaranCore yang berisi data mata pelajaran.
// Jika tidak ada data maka fungsi ini akan mengembalikan error.
func (m *mataPelajaranQuery) SelectMapelById(id string) (*matapelajaran.MataPelajaranCore, error) {
	if m.db == nil {
		// Jika database tidak ada, kembalikan error.
		return nil, errors.New("Nil database")
	}
	if id == "" {
		// Jika id kosong, kembalikan error.
		return nil, errors.New("ID cannot be empty")
	}

	// Buat query untuk mengambil data mata pelajaran berdasarkan id.
	query := `SELECT 
		mp.id,
		mp.nama_pelajaran,
		mp.id_guru,
		g.nama AS nama_guru,
		mp.kelas_id,
		k.kelas AS nama_kelas,
		mp.deskripsi
	FROM mata_pelajaran mp
	LEFT JOIN guru g ON mp.id_guru = g.id
	LEFT JOIN kelas k ON mp.kelas_id = k.id
	WHERE mp.id = $1 AND mp.delete_at IS NULL`

	// Deklarasikan variabel mp dengan tipe MataPelajaranCore untuk menyimpan hasil query
	var mp matapelajaran.MataPelajaranCore

	// Jalankan query untuk mendapatkan data mata pelajaran berdasarkan id, dan pindai hasilnya ke dalam variabel mp
	// Fungsi QueryRow digunakan untuk mengeksekusi query yang mengembalikan satu baris hasil.
	// Kemudian, fungsi Scan digunakan untuk memindai hasil query ke dalam variabel mp.
	err := m.db.QueryRow(context.Background(), query, id).Scan(
		&mp.ID,             // Memindai ID mata pelajaran
		&mp.Nama_Pelajaran, // Memindai nama mata pelajaran
		&mp.ID_Guru,        // Memindai ID guru
		&mp.Guru,           // Memindai nama guru
		&mp.Kelas_ID,       // Memindai ID kelas
		&mp.Nama_Kelas,     // Memindai nama kelas
		&mp.Deskripsi,      // Memindai deskripsi mata pelajaran
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Jika tidak ada data maka kembalikan error.
			return nil, errors.New("ID not found")
		}
		// Jika terjadi error saat query maka log error dan kembalikan.
		log.Printf("QueryRow error: %v", err)
		return nil, fmt.Errorf("select failed: %w", err)
	}
	// Jika data berhasil diambil maka log pesan sukses dan kembalikan data.
	log.Printf("Successfully fetched mata pelajaran with id: %s", id)
	return &mp, nil
}

// UpdateMapel implements matapelajaran.DataMataPelajaranInterface.
// UpdateMapel implements matapelajaran.DataMataPelajaranInterface.
// Fungsi ini digunakan untuk mengupdate data mata pelajaran berdasarkan id.
// Fungsi ini akan mengembalikan error jika terjadi kesalahan.
func (m *mataPelajaranQuery) UpdateMapel(update *matapelajaran.MataPelajaranCore, id string) error {
	// Cek apakah database ada atau tidak.
	if m.db == nil {
		return errors.New("Nil database")
	}
	// Cek apakah id kosong atau tidak.
	if id == "" {
		return errors.New("ID cannot be empty")
	}

	// Buat query untuk mengupdate data mata pelajaran berdasarkan id.
	// Query ini akan mengupdate nama_pelajaran, id_guru, kelas_id, dan deskripsi.
	// Dan akan mengupdate update_at dengan waktu sekarang.
	query := `
	UPDATE mata_pelajaran 
	SET nama_pelajaran = $1,
		id_guru = $2,
		kelas_id = $3,
		deskripsi = $4,
		update_at = CURRENT_TIMESTAMP
	WHERE id = $5 AND delete_at IS NULL;
	`

	// Jalankan query untuk mengupdate data mata pelajaran.
	// Fungsi Exec digunakan untuk mengeksekusi query yang tidak mengembalikan hasil.
	res, err := m.db.Exec(
		context.Background(), query,
		update.Nama_Pelajaran,
		update.ID_Guru,
		update.Kelas_ID,
		update.Deskripsi,
		id,
	)
	if err != nil {
		// Jika terjadi error saat query maka log error dan kembalikan.
		log.Printf("UpdateMapel error exec: %v", err)
		return fmt.Errorf("update failed: %w", err)
	}
	if res.RowsAffected() == 0 {
		// Jika tidak ada baris yang terpengaruh maka log dan kembalikan error.
		log.Printf("UpdateUser: no rows updated for id %s", id)
		return errors.New("update failed: no rows affected")
	}
	// Jika data berhasil diupdate maka log pesan sukses dan kembalikan nil.
	log.Printf("Successfully updated mata_pelajaran with id: %s", id)
	return nil
}

// DeleteMapel mengupdate kolom delete_at dengan waktu sekarang pada data mata pelajaran
// yang dicari berdasarkan id. Jika data berhasil diupdate maka fungsi ini akan
// mengembalikan nil. Jika tidak ada data yang terpengaruh maka fungsi ini akan
// mengembalikan error. Jika terjadi error saat query maka fungsi ini akan log
// error dan kembalikan error.
func (m *mataPelajaranQuery) DeleteMapel(id string) error {
	if m.db == nil {
		// Jika database tidak ada maka kembalikan error.
		return errors.New("Nil database")
	}
	if id == "" {
		// Jika id kosong maka kembalikan error.
		return errors.New("ID tidak boleh kosong")
	}

	// Buat query untuk mengupdate kolom delete_at dengan waktu sekarang.
	// Query ini menggunakan parameter $1 untuk menggantikan nilai id.
	query := "UPDATE mata_pelajaran SET delete_at = NOW() WHERE id = $1 AND delete_at IS NULL"

	// Jalankan query dan simpan hasilnya dalam res.
	// Fungsi Exec digunakan untuk mengeksekusi query yang tidak mengembalikan hasil.
	res, err := m.db.Exec(context.Background(), query, id)
	if err != nil {
		// Jika terjadi error saat query maka log error dan kembalikan.
		log.Printf("DeleteMapel error exec: %v", err)
		return fmt.Errorf("delete failed: %w", err)
	}
	if res.RowsAffected() == 0 {
		// Jika tidak ada baris yang terpengaruh maka log dan kembalikan error.
		log.Printf("DeleteMapel: no rows deleted for id %s", id)
		return errors.New("delete failed: no rows affected")
	}

	// Jika data berhasil diupdate maka log pesan sukses dan kembalikan nil.
	log.Printf("Successfully deleted mata_pelajaran with id: %s", id)
	return nil
}
