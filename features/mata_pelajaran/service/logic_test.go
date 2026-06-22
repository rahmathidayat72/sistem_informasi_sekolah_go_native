package service

import (
	"errors"
	matapelajaran "go_rest_native_sekolah/features/mata_pelajaran"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock untuk DataMataPelajaranInterface
type mockDataMataPelajaran struct {
	mock.Mock
}

func (m *mockDataMataPelajaran) SelectAllMapel() ([]matapelajaran.MataPelajaranCore, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]matapelajaran.MataPelajaranCore), args.Error(1)
}

func (m *mockDataMataPelajaran) SelectMapelById(id string) (*matapelajaran.MataPelajaranCore, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*matapelajaran.MataPelajaranCore), args.Error(1)
}

func (m *mockDataMataPelajaran) InsertMapel(insert *matapelajaran.MataPelajaranCore) error {
	args := m.Called(insert)
	return args.Error(0)
}

func (m *mockDataMataPelajaran) UpdateMapel(insert *matapelajaran.MataPelajaranCore, id string) error {
	args := m.Called(insert, id)
	return args.Error(0)
}

func (m *mockDataMataPelajaran) DeleteMapel(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Test SelectAllMapel
func TestSelectAllMapel(t *testing.T) {
	mockRepo := new(mockDataMataPelajaran)

	t.Run("success get all mapel", func(t *testing.T) {
		expectedMapel := []matapelajaran.MataPelajaranCore{
			{
				ID:             "mapel-001",
				Nama_Pelajaran: "Matematika",
				ID_Guru:        "guru-001",
				Guru:           "Budi Santoso",
				Kelas_ID:       "kelas-001",
				Nama_Kelas:     "10A",
				Deskripsi:      "Pembelajaran Matematika dasar",
				Update_At:      "2024-06-21",
			},
			{
				ID:             "mapel-002",
				Nama_Pelajaran: "Bahasa Indonesia",
				ID_Guru:        "guru-002",
				Guru:           "Siti Nurhaliza",
				Kelas_ID:       "kelas-001",
				Nama_Kelas:     "10A",
				Deskripsi:      "Pembelajaran Bahasa Indonesia",
				Update_At:      "2024-06-21",
			},
		}

		mockRepo.On("SelectAllMapel").Return(expectedMapel, nil).Once()

		svc := &mataPelajaranServiceinterface{mataPelajaranData: mockRepo}
		result, err := svc.SelectAllMapel()

		assert.NoError(t, err)
		assert.Equal(t, expectedMapel, result)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed get all mapel - repository error", func(t *testing.T) {
		mockRepo.On("SelectAllMapel").Return(nil, errors.New("database error")).Once()

		svc := &mataPelajaranServiceinterface{mataPelajaranData: mockRepo}
		result, err := svc.SelectAllMapel()

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed - nil repository", func(t *testing.T) {
		svc := &mataPelajaranServiceinterface{mataPelajaranData: nil}
		result, err := svc.SelectAllMapel()

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// Test SelectMapelById
func TestSelectMapelById(t *testing.T) {
	mockRepo := new(mockDataMataPelajaran)

	t.Run("success get mapel by id", func(t *testing.T) {
		expectedMapel := &matapelajaran.MataPelajaranCore{
			ID:             "mapel-001",
			Nama_Pelajaran: "Matematika",
			ID_Guru:        "guru-001",
			Guru:           "Budi Santoso",
			Kelas_ID:       "kelas-001",
			Nama_Kelas:     "10A",
			Deskripsi:      "Pembelajaran Matematika dasar",
			Update_At:      "2024-06-21",
		}

		mockRepo.On("SelectMapelById", "mapel-001").Return(expectedMapel, nil).Once()

		svc := &mataPelajaranServiceinterface{mataPelajaranData: mockRepo}
		result, err := svc.SelectMapelById("mapel-001")

		assert.NoError(t, err)
		assert.Equal(t, expectedMapel, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed get mapel by id - not found", func(t *testing.T) {
		mockRepo.On("SelectMapelById", "999").Return(nil, pgx.ErrNoRows).Once()

		svc := &mataPelajaranServiceinterface{mataPelajaranData: mockRepo}
		result, err := svc.SelectMapelById("999")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

// Test InsertMapel
func TestInsertMapel(t *testing.T) {
	mockRepo := new(mockDataMataPelajaran)

	t.Run("success insert mapel", func(t *testing.T) {
		newMapel := &matapelajaran.MataPelajaranCore{
			ID:             "mapel-001",
			Nama_Pelajaran: "Matematika",
			ID_Guru:        "guru-001",
			Kelas_ID:       "kelas-001",
			Deskripsi:      "Pembelajaran Matematika dasar",
		}

		mockRepo.On("InsertMapel", newMapel).Return(nil).Once()

		svc := &mataPelajaranServiceinterface{mataPelajaranData: mockRepo}
		err := svc.InsertMapel(newMapel)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed insert mapel - nil repository", func(t *testing.T) {
		newMapel := &matapelajaran.MataPelajaranCore{
			Nama_Pelajaran: "Matematika",
			ID_Guru:        "guru-001",
			Kelas_ID:       "kelas-001",
		}

		svc := &mataPelajaranServiceinterface{mataPelajaranData: nil}
		err := svc.InsertMapel(newMapel)

		assert.Error(t, err)
	})

	t.Run("failed insert mapel - repository error", func(t *testing.T) {
		newMapel := &matapelajaran.MataPelajaranCore{
			Nama_Pelajaran: "Matematika",
			ID_Guru:        "guru-001",
			Kelas_ID:       "kelas-001",
		}

		mockRepo.On("InsertMapel", newMapel).Return(errors.New("insert failed")).Once()

		svc := &mataPelajaranServiceinterface{mataPelajaranData: mockRepo}
		err := svc.InsertMapel(newMapel)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// Test UpdateMapel
func TestUpdateMapel(t *testing.T) {
	mockRepo := new(mockDataMataPelajaran)

	t.Run("success update mapel", func(t *testing.T) {
		existingMapel := &matapelajaran.MataPelajaranCore{
			ID:             "mapel-001",
			Nama_Pelajaran: "Matematika",
			ID_Guru:        "guru-001",
			Kelas_ID:       "kelas-001",
		}

		updatedMapel := &matapelajaran.MataPelajaranCore{
			Nama_Pelajaran: "Matematika Lanjutan",
			ID_Guru:        "guru-002",
			Deskripsi:      "Pembelajaran Matematika lanjutan",
		}

		mockRepo.On("SelectMapelById", "mapel-001").Return(existingMapel, nil).Once()
		mockRepo.On("UpdateMapel", updatedMapel, "mapel-001").Return(nil).Once()

		svc := &mataPelajaranServiceinterface{mataPelajaranData: mockRepo}
		err := svc.UpdateMapel(updatedMapel, "mapel-001")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed update mapel - not found", func(t *testing.T) {
		mockRepo.On("SelectMapelById", "999").Return(nil, pgx.ErrNoRows).Once()

		svc := &mataPelajaranServiceinterface{mataPelajaranData: mockRepo}
		err := svc.UpdateMapel(&matapelajaran.MataPelajaranCore{}, "999")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// Test DeleteMapel
func TestDeleteMapel(t *testing.T) {
	mockRepo := new(mockDataMataPelajaran)

	t.Run("success delete mapel", func(t *testing.T) {
		mockRepo.On("DeleteMapel", "mapel-001").Return(nil).Once()

		svc := &mataPelajaranServiceinterface{mataPelajaranData: mockRepo}
		err := svc.DeleteMapel("mapel-001")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed delete mapel - not found", func(t *testing.T) {
		mockRepo.On("DeleteMapel", "999").Return(errors.New("data not found")).Once()

		svc := &mataPelajaranServiceinterface{mataPelajaranData: mockRepo}
		err := svc.DeleteMapel("999")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
