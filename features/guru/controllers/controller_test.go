package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"go_rest_native_sekolah/features/guru"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock untuk ServiceGuruInterface
type mockServiceGuru struct {
	mock.Mock
}

func (m *mockServiceGuru) GetAllGuru() ([]guru.GuruCore, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]guru.GuruCore), args.Error(1)
}

func (m *mockServiceGuru) InsertGuru(insert *guru.GuruCore) error {
	args := m.Called(insert)
	return args.Error(0)
}

func (m *mockServiceGuru) UpdateGuru(insert *guru.GuruCore, id string) error {
	args := m.Called(insert, id)
	return args.Error(0)
}

func (m *mockServiceGuru) SelectById(id string) (*guru.GuruCore, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*guru.GuruCore), args.Error(1)
}

func (m *mockServiceGuru) DeleteById(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Test GetAllGuru Controller
func TestGetAllGuruController(t *testing.T) {
	mockService := new(mockServiceGuru)

	t.Run("success get all guru", func(t *testing.T) {
		expectedGurus := []guru.GuruCore{
			{
				ID:        "1",
				ID_User:   "user1",
				Nama:      "John Doe",
				Email:     "john@example.com",
				Alamat:    "Jl. Merdeka No. 1",
				Update_At: time.Now(),
			},
		}

		mockService.On("GetAllGuru").Return(expectedGurus, nil).Once()

		controller := NewGuruController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/guru", nil)

		err := controller.Guru(w, r)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("failed get all guru - service error", func(t *testing.T) {
		mockService.On("GetAllGuru").Return(nil, errors.New("service error")).Once()

		controller := NewGuruController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/guru", nil)

		err := controller.Guru(w, r)

		assert.Error(t, err)
		mockService.AssertExpectations(t)
	})
}

// Test InsertGuru Controller
func TestInsertGuruController(t *testing.T) {
	mockService := new(mockServiceGuru)

	t.Run("success insert guru", func(t *testing.T) {
		newGuru := guru.GuruCore{

			Nama:    "John Doe",
			Email:   "john@example.com",
			Alamat:  "Jl. Merdeka No. 1",
		}

		mockService.On("InsertGuru", &newGuru).Return(nil).Once()

		requestBody, _ := json.Marshal(GuruFormatter{
			Nama:   "John Doe",
			Email:  "john@example.com",
			Alamat: "Jl. Merdeka No. 1",
		})

		controller := NewGuruController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/guru/tambah", bytes.NewReader(requestBody))
		r.Header.Set("Content-Type", "application/json")

		err := controller.InsertGuru(w, r)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("failed insert guru - empty required fields", func(t *testing.T) {
		requestBody, _ := json.Marshal(GuruFormatter{
			Nama:   "",
			Email:  "john@example.com",
			Alamat: "Jl. Merdeka No. 1",
		})

		controller := NewGuruController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/guru/tambah", bytes.NewReader(requestBody))
		r.Header.Set("Content-Type", "application/json")

		err := controller.InsertGuru(w, r)

		assert.Error(t, err)
	})
}

// Test GetGuruById Controller
func TestGetGuruByIdController(t *testing.T) {
	mockService := new(mockServiceGuru)

	t.Run("success get guru by id", func(t *testing.T) {
		expectedGuru := &guru.GuruCore{
			ID:        "guru-001",
			Nama:      "John Doe",
			Email:     "john@example.com",
			Alamat:    "Jl. Merdeka No. 1",
			Update_At: time.Now(),
		}

		mockService.On("SelectById", "guru-001").Return(expectedGuru, nil).Once()

		controller := NewGuruController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/guru?id=guru-001", nil)

		err := controller.GetGuruById(w, r)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("failed get guru by id - missing id parameter", func(t *testing.T) {
		controller := NewGuruController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/guru", nil)

		err := controller.GetGuruById(w, r)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("failed get guru by id - not found", func(t *testing.T) {
		mockService.On("SelectById", "deleted-id").Return(nil, errors.New("guru service: Data tidak ditemukan")).Once()

		controller := NewGuruController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/guru/gurubyid?id=deleted-id", nil)

		err := controller.GetGuruById(w, r)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

// Test UpdateGuru Controller
func TestUpdateGuruController(t *testing.T) {
	mockService := new(mockServiceGuru)

	t.Run("success update guru", func(t *testing.T) {
		existingGuru := &guru.GuruCore{
			ID:     "guru-001",
			Nama:   "John Doe",
			Email:  "john@example.com",
			Alamat: "Jl. Merdeka No. 1",
		}

		updatedGuru := guru.GuruCore{
			Nama:   "John Updated",
			Email:  "john.updated@example.com",
			Alamat: "Jl. Merdeka No. 2",
		}

		mockService.On("SelectById", "guru-001").Return(existingGuru, nil).Once()
		mockService.On("UpdateGuru", &updatedGuru, "guru-001").Return(nil).Once()
		mockService.On("SelectById", "guru-001").Return(&updatedGuru, nil).Once()

		requestBody, _ := json.Marshal(GuruFormatter{
			Nama:   "John Updated",
			Email:  "john.updated@example.com",
			Alamat: "Jl. Merdeka No. 2",
		})

		controller := NewGuruController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/guru?id=guru-001", bytes.NewReader(requestBody))
		r.Header.Set("Content-Type", "application/json")

		err := controller.UpdateGuru(w, r)

		assert.NoError(t, err)
	})

	t.Run("failed update guru - missing id parameter", func(t *testing.T) {
		controller := NewGuruController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/guru", nil)

		err := controller.UpdateGuru(w, r)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// Test DeleteGuru Controller
func TestDeleteGuruController(t *testing.T) {
	mockService := new(mockServiceGuru)

	t.Run("success delete guru", func(t *testing.T) {
		mockService.On("DeleteById", "guru-001").Return(nil).Once()

		controller := NewGuruController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/guru?id=guru-001", nil)

		err := controller.DeleteGuru(w, r)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("failed delete guru - missing id parameter", func(t *testing.T) {
		controller := NewGuruController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/guru", nil)

		err := controller.DeleteGuru(w, r)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
