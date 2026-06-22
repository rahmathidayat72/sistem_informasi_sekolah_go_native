package service

import (
	"errors"
	"go_rest_native_sekolah/features/guru"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock untuk DataGuruInterface
type mockDataGuru struct {
	mock.Mock
}

func (m *mockDataGuru) SelectAllGuru() ([]guru.GuruCore, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]guru.GuruCore), args.Error(1)
}

func (m *mockDataGuru) InsertGuru(insert *guru.GuruCore) error {
	args := m.Called(insert)
	return args.Error(0)
}

func (m *mockDataGuru) Update(insert *guru.GuruCore, id string) error {
	args := m.Called(insert, id)
	return args.Error(0)
}

func (m *mockDataGuru) SelectById(id string) (*guru.GuruCore, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*guru.GuruCore), args.Error(1)
}

func (m *mockDataGuru) DeleteById(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Test GetAllGuru
func TestGetAllGuru(t *testing.T) {
	mockRepo := new(mockDataGuru)

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
			{
				ID:        "2",
				ID_User:   "user2",
				Nama:      "Jane Smith",
				Email:     "jane@example.com",
				Alamat:    "Jl. Sudirman No. 2",
				Update_At: time.Now(),
			},
		}

		mockRepo.On("SelectAllGuru").Return(expectedGurus, nil).Once()

		svc := &guruService{guruData: mockRepo, db: nil}
		result, err := svc.GetAllGuru()

		assert.NoError(t, err)
		assert.Equal(t, expectedGurus, result)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed get all guru - repository error", func(t *testing.T) {
		mockRepo.On("SelectAllGuru").Return(nil, errors.New("database error")).Once()

		svc := &guruService{guruData: mockRepo, db: nil}
		result, err := svc.GetAllGuru()

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "gagal mengambil data")
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed - nil repository", func(t *testing.T) {
		svc := &guruService{guruData: nil, db: nil}
		result, err := svc.GetAllGuru()

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Nil repository")
	})
}

// Test InsertGuru
func TestInsertGuru(t *testing.T) {
	mockRepo := new(mockDataGuru)

	t.Run("failed insert guru - invalid email", func(t *testing.T) {
		invalidGuru := &guru.GuruCore{
			ID:      "guru-001",
			Nama:    "John Doe",
			Email:   "invalid-email",
			Alamat:  "Jl. Merdeka No. 1",
		}

		svc := &guruService{guruData: mockRepo, db: nil}
		err := svc.InsertGuru(invalidGuru)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email tidak valid")
	})

	t.Run("failed insert guru - empty required fields (nama)", func(t *testing.T) {
		invalidGuru := &guru.GuruCore{
			ID:      "guru-001",
			Nama:    "",
			Email:   "john@example.com",
			Alamat:  "Jl. Merdeka No. 1",
		}

		svc := &guruService{guruData: mockRepo, db: nil}
		err := svc.InsertGuru(invalidGuru)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "harus diisi")
	})

	t.Run("failed insert guru - empty required fields (email)", func(t *testing.T) {
		invalidGuru := &guru.GuruCore{
			ID:      "guru-001",
			Nama:    "John Doe",
			Email:   "",
			Alamat:  "Jl. Merdeka No. 1",
		}

		svc := &guruService{guruData: mockRepo, db: nil}
		err := svc.InsertGuru(invalidGuru)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "harus diisi")
	})

	t.Run("failed insert guru - empty required fields (alamat)", func(t *testing.T) {
		invalidGuru := &guru.GuruCore{
			ID:      "guru-001",
			Nama:    "John Doe",
			Email:   "john@example.com",
			Alamat:  "",
		}

		svc := &guruService{guruData: mockRepo, db: nil}
		err := svc.InsertGuru(invalidGuru)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "harus diisi")
	})

	t.Run("failed insert guru - nil repository", func(t *testing.T) {
		newGuru := &guru.GuruCore{
			Nama:    "John Doe",
			Email:   "john@example.com",
			Alamat:  "Jl. Merdeka No. 1",
		}

		svc := &guruService{guruData: nil, db: nil}
		err := svc.InsertGuru(newGuru)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Repository kosong")
	})
}

// Test SelectById
func TestSelectGuruById(t *testing.T) {
	mockRepo := new(mockDataGuru)

	t.Run("success get guru by id", func(t *testing.T) {
		expectedGuru := &guru.GuruCore{
			ID:        "1",
			ID_User:   "user1",
			Nama:      "John Doe",
			Email:     "john@example.com",
			Alamat:    "Jl. Merdeka No. 1",
			Update_At: time.Now(),
		}

		mockRepo.On("SelectById", "1").Return(expectedGuru, nil).Once()

		svc := &guruService{guruData: mockRepo, db: nil}
		result, err := svc.SelectById("1")

		assert.NoError(t, err)
		assert.Equal(t, expectedGuru, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed get guru by id - not found", func(t *testing.T) {
		mockRepo.On("SelectById", "999").Return(nil, pgx.ErrNoRows).Once()

		svc := &guruService{guruData: mockRepo, db: nil}
		result, err := svc.SelectById("999")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

// Test UpdateGuru
func TestUpdateGuru(t *testing.T) {
	mockRepo := new(mockDataGuru)

	t.Run("success update guru", func(t *testing.T) {
		existingGuru := &guru.GuruCore{
			ID:      "1",
			Nama:    "John Doe",
			Email:   "john@example.com",
			Alamat:  "Jl. Merdeka No. 1",
		}

		updatedGuru := &guru.GuruCore{
			Nama:   "John Updated",
			Email:  "john.updated@example.com",
			Alamat: "Jl. Merdeka No. 2",
		}

		mockRepo.On("SelectById", "1").Return(existingGuru, nil).Once()
		mockRepo.On("Update", updatedGuru, "1").Return(nil).Once()

		svc := &guruService{guruData: mockRepo, db: nil}
		err := svc.UpdateGuru(updatedGuru, "1")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed update guru - not found", func(t *testing.T) {
		mockRepo.On("SelectById", "999").Return(nil, pgx.ErrNoRows).Once()

		svc := &guruService{guruData: mockRepo, db: nil}
		err := svc.UpdateGuru(&guru.GuruCore{}, "999")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tidak ditemukan")
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed update guru - empty id", func(t *testing.T) {
		svc := &guruService{guruData: mockRepo, db: nil}
		err := svc.UpdateGuru(&guru.GuruCore{}, "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "id harus diisi")
	})
}

// Test DeleteById
func TestDeleteGuruById(t *testing.T) {
	mockRepo := new(mockDataGuru)

	t.Run("success delete guru", func(t *testing.T) {
		mockRepo.On("DeleteById", "1").Return(nil).Once()

		svc := &guruService{guruData: mockRepo, db: nil}
		err := svc.DeleteById("1")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed delete guru - not found", func(t *testing.T) {
		mockRepo.On("DeleteById", "999").Return(errors.New("data not found")).Once()

		svc := &guruService{guruData: mockRepo, db: nil}
		err := svc.DeleteById("999")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// Test Panic when nil repository
func TestNewServiceGuruPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic but got none")
		}
	}()

	NewServiceGuru(nil, &pgxpool.Pool{})
}
