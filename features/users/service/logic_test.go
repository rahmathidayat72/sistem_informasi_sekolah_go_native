package service

import (
	"errors"
	"go_rest_native_sekolah/features/users"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock untuk DataUserInterface
type mockDataUser struct {
	mock.Mock
}

func (m *mockDataUser) SelectAllUser() ([]users.UserCore, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]users.UserCore), args.Error(1)
}

func (m *mockDataUser) SelectUserById(id string) (*users.UserCore, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.UserCore), args.Error(1)
}

func (m *mockDataUser) InsertUser(input *users.UserCore) error {
	args := m.Called(input)
	return args.Error(0)
}

func (m *mockDataUser) UpdateUser(input *users.UserCore, id string) error {
	args := m.Called(input, id)
	return args.Error(0)
}

func (m *mockDataUser) DeleteUserById(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Test SelectAllUser
func TestSelectAllUser(t *testing.T) {
	mockRepo := new(mockDataUser)

	t.Run("success get all users", func(t *testing.T) {
		expectedUsers := []users.UserCore{
			{
				ID:        "user-001",
				Username:  "john_doe",
				Email:     "john@example.com",
				Password:  "hashed_password",
				Role:      "admin",
				Update_At: time.Now(),
			},
			{
				ID:        "user-002",
				Username:  "jane_smith",
				Email:     "jane@example.com",
				Password:  "hashed_password",
				Role:      "user",
				Update_At: time.Now(),
			},
		}

		mockRepo.On("SelectAllUser").Return(expectedUsers, nil).Once()

		svc := &userService{userData: mockRepo}
		result, err := svc.SelectAllUser()

		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, result)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed get all users - repository error", func(t *testing.T) {
		mockRepo.On("SelectAllUser").Return(nil, errors.New("database error")).Once()

		svc := &userService{userData: mockRepo}
		result, err := svc.SelectAllUser()

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed - nil repository", func(t *testing.T) {
		svc := &userService{userData: nil}
		result, err := svc.SelectAllUser()

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// Test SelectUserById
func TestSelectUserById(t *testing.T) {
	mockRepo := new(mockDataUser)

	t.Run("success get user by id", func(t *testing.T) {
		expectedUser := &users.UserCore{
			ID:        "user-001",
			Username:  "john_doe",
			Email:     "john@example.com",
			Password:  "hashed_password",
			Role:      "admin",
			Update_At: time.Now(),
		}

		mockRepo.On("SelectUserById", "user-001").Return(expectedUser, nil).Once()

		svc := &userService{userData: mockRepo}
		result, err := svc.SelectUserById("user-001")

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed get user by id - not found", func(t *testing.T) {
		mockRepo.On("SelectUserById", "999").Return(nil, pgx.ErrNoRows).Once()

		svc := &userService{userData: mockRepo}
		result, err := svc.SelectUserById("999")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

// Test InsertUser
func TestInsertUser(t *testing.T) {
	mockRepo := new(mockDataUser)

	t.Run("success insert user", func(t *testing.T) {
		newUser := &users.UserCore{
			ID:       "user-001",
			Username: "john_doe",
			Email:    "john@example.com",
			Password: "hashed_password",
			Role:     "admin",
		}

		mockRepo.On("InsertUser", newUser).Return(nil).Once()

		svc := &userService{userData: mockRepo}
		err := svc.InsertUser(newUser)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed insert user - nil repository", func(t *testing.T) {
		newUser := &users.UserCore{
			Username: "john_doe",
			Email:    "john@example.com",
			Password: "hashed_password",
			Role:     "admin",
		}

		svc := &userService{userData: nil}
		err := svc.InsertUser(newUser)

		assert.Error(t, err)
	})

	t.Run("failed insert user - repository error", func(t *testing.T) {
		newUser := &users.UserCore{
			Username: "john_doe",
			Email:    "john@example.com",
			Password: "hashed_password",
			Role:     "admin",
		}

		mockRepo.On("InsertUser", newUser).Return(errors.New("insert failed")).Once()

		svc := &userService{userData: mockRepo}
		err := svc.InsertUser(newUser)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// Test UpdateUser
func TestUpdateUser(t *testing.T) {
	mockRepo := new(mockDataUser)

	t.Run("success update user", func(t *testing.T) {
		existingUser := &users.UserCore{
			ID:        "user-001",
			Username:  "john_doe",
			Email:     "john@example.com",
			Password:  "hashed_password",
			Role:      "admin",
			Update_At: time.Now(),
		}

		updatedUser := &users.UserCore{
			Username: "john_updated",
			Email:    "john.updated@example.com",
			Password: "new_hashed_password",
			Role:     "user",
		}

		mockRepo.On("SelectUserById", "user-001").Return(existingUser, nil).Once()
		mockRepo.On("UpdateUser", updatedUser, "user-001").Return(nil).Once()

		svc := &userService{userData: mockRepo}
		err := svc.UpdateUser(updatedUser, "user-001")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed update user - not found", func(t *testing.T) {
		mockRepo.On("SelectUserById", "999").Return(nil, pgx.ErrNoRows).Once()

		svc := &userService{userData: mockRepo}
		err := svc.UpdateUser(&users.UserCore{}, "999")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// Test DeleteUserById
func TestDeleteUserById(t *testing.T) {
	mockRepo := new(mockDataUser)

	t.Run("success delete user", func(t *testing.T) {
		mockRepo.On("DeleteUserById", "user-001").Return(nil).Once()

		svc := &userService{userData: mockRepo}
		err := svc.DeleteUserById("user-001")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed delete user - not found", func(t *testing.T) {
		mockRepo.On("DeleteUserById", "999").Return(errors.New("data not found")).Once()

		svc := &userService{userData: mockRepo}
		err := svc.DeleteUserById("999")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
