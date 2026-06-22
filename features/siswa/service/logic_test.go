package service

import (
	"errors"
	"go_rest_native_sekolah/features/siswa"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock untuk DataSiswaInterface
type mockDataSiswa struct {
	mock.Mock
}

func (m *mockDataSiswa) SelectAllSiswa() ([]siswa.SiswaCore, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]siswa.SiswaCore), args.Error(1)
}

func (m *mockDataSiswa) InsertSiswa(insert *siswa.SiswaCore) error {
	args := m.Called(insert)
	return args.Error(0)
}

func (m *mockDataSiswa) Update(insert *siswa.SiswaCore, id string) error {
	args := m.Called(insert, id)
	return args.Error(0)
}

func (m *mockDataSiswa) SelectById(id string) (*siswa.SiswaCore, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*siswa.SiswaCore), args.Error(1)
}

func (m *mockDataSiswa) DeleteById(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Test SelectAllSiswa
func TestSelectAllSiswa(t *testing.T) {
	mockRepo := new(mockDataSiswa)

	t.Run("success get all siswa", func(t *testing.T) {
		expectedSiswa := []siswa.SiswaCore{
			{
				ID:         "siswa-001",
				Nama:       "Ahmad Rauf",
				Kelas_ID:   "kelas-001",
				Nama_Kelas: "10A",
				Email:      "ahmad@example.com",
				Alamat:     "Jl. Gatot Subroto No. 1",
				Update_At:  time.Now(),
			},
			{
				ID:         "siswa-002",
				Nama:       "Siti Nurhaliza",
				Kelas_ID:   "kelas-001",
				Nama_Kelas: "10A",
				Email:      "siti@example.com",
				Alamat:     "Jl. Ahmad Yani No. 2",
				Update_At:  time.Now(),
			},
		}

		mockRepo.On("SelectAllSiswa").Return(expectedSiswa, nil).Once()

		svc := &siswaService{siswaData: mockRepo}
		result, err := svc.SelectAllSiswa()

		assert.NoError(t, err)
		assert.Equal(t, expectedSiswa, result)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed get all siswa - repository error", func(t *testing.T) {
		mockRepo.On("SelectAllSiswa").Return(nil, errors.New("database error")).Once()

		svc := &siswaService{siswaData: mockRepo}
		result, err := svc.SelectAllSiswa()

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed - nil repository", func(t *testing.T) {
		svc := &siswaService{siswaData: nil}
		result, err := svc.SelectAllSiswa()

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// Test InsertSiswa
func TestInsertSiswa(t *testing.T) {
	mockRepo := new(mockDataSiswa)

	t.Run("success insert siswa", func(t *testing.T) {
		newSiswa := &siswa.SiswaCore{
			ID:       "siswa-001",
			Nama:     "Ahmad Rauf",
			Kelas_ID: "kelas-001",
			Email:    "ahmad@example.com",
			Alamat:   "Jl. Gatot Subroto No. 1",
		}

		mockRepo.On("InsertSiswa", newSiswa).Return(nil).Once()

		svc := &siswaService{siswaData: mockRepo}
		err := svc.InsertSiswa(newSiswa)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed insert siswa - invalid email", func(t *testing.T) {
		invalidSiswa := &siswa.SiswaCore{
			ID:       "siswa-001",
			Nama:     "Ahmad Rauf",
			Kelas_ID: "kelas-001",
			Email:    "invalid-email",
			Alamat:   "Jl. Gatot Subroto No. 1",
		}

		svc := &siswaService{siswaData: mockRepo}
		err := svc.InsertSiswa(invalidSiswa)

		assert.Error(t, err)
	})

	t.Run("failed insert siswa - nil repository", func(t *testing.T) {
		newSiswa := &siswa.SiswaCore{
			Nama:     "Ahmad Rauf",
			Kelas_ID: "kelas-001",
			Email:    "ahmad@example.com",
			Alamat:   "Jl. Gatot Subroto No. 1",
		}

		svc := &siswaService{siswaData: nil}
		err := svc.InsertSiswa(newSiswa)

		assert.Error(t, err)
	})
}

// Test SelectSiswaById
func TestSelectSiswaById(t *testing.T) {
	mockRepo := new(mockDataSiswa)

	t.Run("success get siswa by id", func(t *testing.T) {
		expectedSiswa := &siswa.SiswaCore{
			ID:         "siswa-001",
			Nama:       "Ahmad Rauf",
			Kelas_ID:   "kelas-001",
			Nama_Kelas: "10A",
			Email:      "ahmad@example.com",
			Alamat:     "Jl. Gatot Subroto No. 1",
			Update_At:  time.Now(),
		}

		mockRepo.On("SelectById", "siswa-001").Return(expectedSiswa, nil).Once()

		svc := &siswaService{siswaData: mockRepo}
		result, err := svc.SelectById("siswa-001")

		assert.NoError(t, err)
		assert.Equal(t, expectedSiswa, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed get siswa by id - not found", func(t *testing.T) {
		mockRepo.On("SelectById", "999").Return(nil, pgx.ErrNoRows).Once()

		svc := &siswaService{siswaData: mockRepo}
		result, err := svc.SelectById("999")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

// Test UpdateSiswa
func TestUpdateSiswa(t *testing.T) {
	mockRepo := new(mockDataSiswa)

	t.Run("success update siswa", func(t *testing.T) {
		existingSiswa := &siswa.SiswaCore{
			ID:       "siswa-001",
			Nama:     "Ahmad Rauf",
			Kelas_ID: "kelas-001",
			Email:    "ahmad@example.com",
			Alamat:   "Jl. Gatot Subroto No. 1",
		}

		updatedSiswa := &siswa.SiswaCore{
			Nama:   "Ahmad Rauf Updated",
			Email:  "ahmad.updated@example.com",
			Alamat: "Jl. Gatot Subroto No. 2",
		}

		mockRepo.On("SelectById", "siswa-001").Return(existingSiswa, nil).Once()
		mockRepo.On("Update", updatedSiswa, "siswa-001").Return(nil).Once()

		svc := &siswaService{siswaData: mockRepo}
		err := svc.Update(updatedSiswa, "siswa-001")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed update siswa - not found", func(t *testing.T) {
		mockRepo.On("SelectById", "999").Return(nil, pgx.ErrNoRows).Once()

		svc := &siswaService{siswaData: mockRepo}
		err := svc.Update(&siswa.SiswaCore{}, "999")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// Test DeleteSiswaById
func TestDeleteSiswaById(t *testing.T) {
	mockRepo := new(mockDataSiswa)

	t.Run("success delete siswa", func(t *testing.T) {
		mockRepo.On("DeleteById", "siswa-001").Return(nil).Once()

		svc := &siswaService{siswaData: mockRepo}
		err := svc.DeleteById("siswa-001")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed delete siswa - not found", func(t *testing.T) {
		mockRepo.On("DeleteById", "999").Return(errors.New("data not found")).Once()

		svc := &siswaService{siswaData: mockRepo}
		err := svc.DeleteById("999")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
