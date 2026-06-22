package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"go_rest_native_sekolah/features/users"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock untuk ServiceUserInterface
type mockServiceUser struct {
	mock.Mock
}

func (m *mockServiceUser) SelectAllUser() ([]users.UserCore, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]users.UserCore), args.Error(1)
}

func (m *mockServiceUser) SelectUserById(id string) (*users.UserCore, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.UserCore), args.Error(1)
}

func (m *mockServiceUser) InsertUser(input *users.UserCore) error {
	args := m.Called(input)
	return args.Error(0)
}

func (m *mockServiceUser) UpdateUser(input *users.UserCore, id string) error {
	args := m.Called(input, id)
	return args.Error(0)
}

func (m *mockServiceUser) DeleteUserById(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Test Users Controller
func TestUsersController(t *testing.T) {
	mockService := new(mockServiceUser)

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
		}

		mockService.On("SelectAllUser").Return(expectedUsers, nil).Once()

		controller := NewUsesController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/users", nil)

		err := controller.Users(w, r)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("failed get all users - no data found", func(t *testing.T) {
		mockService.On("SelectAllUser").Return([]users.UserCore{}, errors.New("no data")).Once()

		controller := NewUsesController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/users", nil)

		err := controller.Users(w, r)

		assert.Error(t, err)
	})
}

// Test InsertUser Controller
func TestInsertUserController(t *testing.T) {
	mockService := new(mockServiceUser)

	t.Run("success insert user", func(t *testing.T) {
		newUser := users.UserCore{
			Username: "john_doe",
			Email:    "john@example.com",
			Password: "password123",
			Role:     "user",
		}

		mockService.On("InsertUser", &newUser).Return(nil).Once()

		requestBody, _ := json.Marshal(UserFormatter{
			Username: "john_doe",
			Email:    "john@example.com",
			Password: "password123",
			Role:     "user",
		})

		controller := NewUsesController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/users/tambah", bytes.NewReader(requestBody))
		r.Header.Set("Content-Type", "application/json")

		err := controller.InsertUser(w, r)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("failed insert user - empty required fields", func(t *testing.T) {
		requestBody, _ := json.Marshal(UserFormatter{
			Username: "",
			Email:    "john@example.com",
			Password: "password123",
			Role:     "user",
		})

		controller := NewUsesController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/users/tambah", bytes.NewReader(requestBody))
		r.Header.Set("Content-Type", "application/json")

		err := controller.InsertUser(w, r)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("failed insert user - invalid JSON", func(t *testing.T) {
		controller := NewUsesController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/users/tambah", bytes.NewReader([]byte("invalid json")))
		r.Header.Set("Content-Type", "application/json")

		err := controller.InsertUser(w, r)

		assert.Error(t, err)
	})
}

// Test GetUserById Controller
func TestGetUserByIdController(t *testing.T) {
	mockService := new(mockServiceUser)

	t.Run("success get user by id", func(t *testing.T) {
		expectedUser := &users.UserCore{
			ID:        "user-001",
			Username:  "john_doe",
			Email:     "john@example.com",
			Password:  "hashed_password",
			Role:      "admin",
			Update_At: time.Now(),
		}

		mockService.On("SelectUserById", "user-001").Return(expectedUser, nil).Once()

		controller := NewUsesController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/users?id=user-001", nil)

		err := controller.GetUserById(w, r)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("failed get user by id - missing id parameter", func(t *testing.T) {
		controller := NewUsesController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/users", nil)

		err := controller.GetUserById(w, r)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("failed get user by id - not found", func(t *testing.T) {
		mockService.On("SelectUserById", "999").Return(nil, errors.New("not found")).Once()

		controller := NewUsesController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/users?id=999", nil)

		err := controller.GetUserById(w, r)

		assert.Error(t, err)
	})
}

// Test UpdateUser Controller
func TestUpdateUserController(t *testing.T) {
	mockService := new(mockServiceUser)

	t.Run("success update user", func(t *testing.T) {
		existingUser := &users.UserCore{
			ID:        "user-001",
			Username:  "john_doe",
			Email:     "john@example.com",
			Password:  "hashed_password",
			Role:      "admin",
			Update_At: time.Now(),
		}

		updatedUser := users.UserCore{
			Username: "john_updated",
			Email:    "john.updated@example.com",
			Password: "new_password",
			Role:     "user",
		}

		mockService.On("SelectUserById", "user-001").Return(existingUser, nil).Once()
		mockService.On("UpdateUser", &updatedUser, "user-001").Return(nil).Once()
		mockService.On("SelectUserById", "user-001").Return(&updatedUser, nil).Once()

		requestBody, _ := json.Marshal(UserFormatter{
			Username: "john_updated",
			Email:    "john.updated@example.com",
			Password: "new_password",
			Role:     "user",
		})

		controller := NewUsesController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/users?id=user-001", bytes.NewReader(requestBody))
		r.Header.Set("Content-Type", "application/json")

		err := controller.UpdateUser(w, r)

		assert.NoError(t, err)
	})

	t.Run("failed update user - missing id parameter", func(t *testing.T) {
		controller := NewUsesController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/users", nil)

		err := controller.UpdateUser(w, r)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// Test DeleteUser Controller
func TestDeleteUserController(t *testing.T) {
	mockService := new(mockServiceUser)

	t.Run("success delete user", func(t *testing.T) {
		mockService.On("DeleteUserById", "user-001").Return(nil).Once()

		controller := NewUsesController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/users?id=user-001", nil)

		err := controller.DeleteUser(w, r)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("failed delete user - missing id parameter", func(t *testing.T) {
		controller := NewUsesController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/users", nil)

		err := controller.DeleteUser(w, r)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("failed delete user - service error", func(t *testing.T) {
		mockService.On("DeleteUserById", "user-001").Return(errors.New("delete failed")).Once()

		controller := NewUsesController(mockService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/users?id=user-001", nil)

		err := controller.DeleteUser(w, r)

		assert.Error(t, err)
	})
}
