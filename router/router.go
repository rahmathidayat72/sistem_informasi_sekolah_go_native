package router

import (
	authcontroller "go_rest_native_sekolah/features/auth/controllers"
	authmodels "go_rest_native_sekolah/features/auth/model"
	serviceauth "go_rest_native_sekolah/features/auth/service"
	gurucontroller "go_rest_native_sekolah/features/guru/controllers"
	gurumodels "go_rest_native_sekolah/features/guru/model"
	"go_rest_native_sekolah/features/guru/service"
	kelascontroller "go_rest_native_sekolah/features/kelas/controllers"
	kelasmodels "go_rest_native_sekolah/features/kelas/model"
	servicekelas "go_rest_native_sekolah/features/kelas/service"
	mapelcontroller "go_rest_native_sekolah/features/mata_pelajaran/controllers"
	mapelsmodels "go_rest_native_sekolah/features/mata_pelajaran/model"
	servicemapel "go_rest_native_sekolah/features/mata_pelajaran/service"
	siswacontroller "go_rest_native_sekolah/features/siswa/controllers"
	siswamodels "go_rest_native_sekolah/features/siswa/model"
	servicesiswa "go_rest_native_sekolah/features/siswa/service"
	userscontroller "go_rest_native_sekolah/features/users/controllers"
	usersmodels "go_rest_native_sekolah/features/users/model"
	serviceuser "go_rest_native_sekolah/features/users/service"

	"go_rest_native_sekolah/helper"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// InitRouter digunakan untuk menginisialisasi router.
// Fungsi ini akan menginisialisasi router untuk fitur auth, guru, users, dan kelas.
func InitRouter(db *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()

	// Pasang semua route
	// Endpoint /login digunakan untuk mengotentikasi user
	loginRouter(mux, db)
	// Endpoint /guru digunakan untuk mengelola data guru
	guruRouter(mux, db)
	// Endpoint /users digunakan untuk mengelola data user
	usersRouter(mux, db)
	// Endpoint /kelas digunakan untuk mengelola data kelas
	kelasRouter(mux, db)
	// Endpoint /siswa digunakan untuk mengelola data siswa
	siswaRouter(mux, db)
	// Endpoint /mapel digunakan untuk mengelola data mata pelajaran
	mataPelajaranRouter(mux, db)

	// Bungkus mux dengan middleware logging
	// Middleware logging digunakan untuk mencatat setiap request yang diterima oleh server
	return helper.LoggingMiddleware(mux, db)
}

// loginRouter digunakan untuk menginisialisasi router untuk fitur auth.
// Fungsi ini akan menginisialisasi router untuk endpoint /login yang digunakan untuk mengotentikasi user.
// Endpoint /login akan menerima request dengan method POST dan mengembalikan response JSON.
func loginRouter(mux *http.ServeMux, db *pgxpool.Pool) {
	// Inisialisasi repository
	authRepo := authmodels.NewAuthData(db)
	// Inisialisasi service
	authService := serviceauth.NewServiceAuth(authRepo)
	// Inisialisasi controller
	authController := authcontroller.NewAutController(authService)
	// Membuat handler untuk endpoint /login
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// Cek apakah request menggunakan method POST
		if r.Method == http.MethodPost {
			err := authController.Auth(w, r)
			if err != nil {
				// Jika terjadi error maka akan mengembalikan response JSON dengan status Internal Server Error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			// Jika request tidak menggunakan method POST maka akan mengembalikan response JSON dengan status Method Not Allowed
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})
}

func guruRouter(mux *http.ServeMux, db *pgxpool.Pool) {
	// Inisialisasi repository
	guruRepo := gurumodels.NewDataGuru(db)

	// Inisialisasi service
	guruService := service.NewServiceGuru(guruRepo, db)

	// Inisialisasi controller
	guruController := gurucontroller.NewGuruController(guruService)

	// Endpoint GET untuk menampilkan data guru
	mux.HandleFunc("/guru", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			err := guruController.Guru(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))

	// Endpoint POST untuk menambah data guru
	mux.HandleFunc("/guru/tambah", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := guruController.InsertGuru(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))

	mux.HandleFunc("/guru/update", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			// Langsung jalankan fungsi UpdateGuru milik controller
			err := guruController.UpdateGuru(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))
	mux.HandleFunc("/guru/gurubyid", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodPut {
			// Langsung jalankan fungsi GuruById milik controller
			err := guruController.GetGuruById(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))

	mux.HandleFunc("/guru/deleted", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete || r.Method == http.MethodPut {
			// Langsung jalankan fungsi DeletedById milik controller
			err := guruController.DeleteGuru(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))
}

func usersRouter(mux *http.ServeMux, db *pgxpool.Pool) {
	usersRepo := usersmodels.NewUserData(db)
	usersService := serviceuser.NewServiceUser(usersRepo, db)
	usersController := userscontroller.NewUsesController(usersService)

	mux.HandleFunc("/users", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			err := usersController.Users(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))

	mux.HandleFunc("/users/tambah", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := usersController.InsertUser(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	mux.HandleFunc("/users/userbyid", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodPut {
			// Langsung jalankan fungsi GuruById milik controller
			err := usersController.GetUserById(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))

	mux.HandleFunc("/users/update", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			// Langsung jalankan fungsi UpdateGuru milik controller
			err := usersController.UpdateUser(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))

	mux.HandleFunc("/users/deleted", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete || r.Method == http.MethodPut {
			// Langsung jalankan fungsi DeletedById milik controller
			err := usersController.DeleteUser(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))
}

func kelasRouter(mux *http.ServeMux, db *pgxpool.Pool) {
	kelasRepo := kelasmodels.NewDataKelas(db)
	kelasService := servicekelas.NewServiceKelas(kelasRepo)
	kelasController := kelascontroller.NewKelasController(kelasService)

	mux.HandleFunc("/kelas/tambah", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := kelasController.Insert(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}

	}))

	mux.HandleFunc("/kelas", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			err := kelasController.Kelas(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))

	mux.HandleFunc("/kelas/kelasbyid", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Langsung jalankan fungsi GuruById milik controller
			err := kelasController.GetKelasById(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))

	mux.HandleFunc("/kelas/update", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			// Langsung jalankan fungsi UpdateGuru milik controller
			err := kelasController.UpdateKelas(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))

	mux.HandleFunc("/kelas/deleted", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete || r.Method == http.MethodPut {
			// Langsung jalankan fungsi DeletedById milik controller
			err := kelasController.DeleteKelas(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))
}

func siswaRouter(mux *http.ServeMux, db *pgxpool.Pool) {
	{
		siswaRepo := siswamodels.NewSiswaData(db)
		siswaService := servicesiswa.NewServiceSiswa(siswaRepo)
		siswaController := siswacontroller.NewSiswaController(siswaService)

		mux.HandleFunc("/siswa/tambah", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				err := siswaController.InsertSiswa(w, r)
				if err != nil {
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			} else {
				helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		}))

		mux.HandleFunc("/siswa", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				err := siswaController.Siswa(w, r)
				if err != nil {
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			} else {
				helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		}))

		mux.HandleFunc("/siswa/siswabyid", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				// Langsung jalankan fungsi GuruById milik controller
				err := siswaController.GetSiswaById(w, r)
				if err != nil {
					// Tampilkan response JSON error
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			} else {
				helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		}))

		mux.HandleFunc("/siswa/update", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost || r.Method == http.MethodPut {
				// Langsung jalankan fungsi UpdateGuru milik controller
				err := siswaController.UpdateSiswa(w, r)
				if err != nil {
					// Tampilkan response JSON error
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			} else {
				helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		}))

		mux.HandleFunc("/siswa/deleted", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete || r.Method == http.MethodPut {
				// Langsung jalankan fungsi DeletedById milik controller
				err := siswaController.DeleteSiswa(w, r)
				if err != nil {
					// Tampilkan response JSON error
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			} else {
				helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		}))

	}
}

func mataPelajaranRouter(mux *http.ServeMux, db *pgxpool.Pool) {
	{
		mataPelajaranRepo := mapelsmodels.NewDataMataPelajaran(db)
		mataPelajaranService := servicemapel.NewMataPelajaranService(mataPelajaranRepo)
		mataPelajaranController := mapelcontroller.NewMataPelajaranController(mataPelajaranService)

		mux.HandleFunc("/mapel/tambah", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				err := mataPelajaranController.InsertMapel(w, r)
				if err != nil {
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			} else {
				helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		}))

		mux.HandleFunc("/mapel", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				err := mataPelajaranController.Mapel(w, r)
				if err != nil {
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			} else {
				helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		}))

		mux.HandleFunc("/mapel/mapelbyid", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				// Langsung jalankan fungsi GuruById milik controller
				err := mataPelajaranController.GetMapelById(w, r)
				if err != nil {
					// Tampilkan response JSON error
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			} else {
				helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		}))

		mux.HandleFunc("/mapel/update", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost || r.Method == http.MethodPut {
				// Langsung jalankan fungsi UpdateGuru milik controller
				err := mataPelajaranController.UpdateMapel(w, r)
				if err != nil {
					// Tampilkan response JSON error
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			} else {
				helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		}))

		mux.HandleFunc("/mapel/deleted", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete || r.Method == http.MethodPut {
				// Langsung jalankan fungsi DeletedById milik controller
				err := mataPelajaranController.DeleteMapel(w, r)
				if err != nil {
					// Tampilkan response JSON error
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			} else {
				helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		}))
	}
}
