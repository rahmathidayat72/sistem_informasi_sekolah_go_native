package service_test

import (
	"errors"
	"go_rest_native_sekolah/features/auth"
	"go_rest_native_sekolah/features/auth/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock untuk DataAuthInterface
type mockDataAuth struct {
	mock.Mock
}

func (m *mockDataAuth) Login(email, password string) (auth.UserCore, error) {
	args := m.Called(email, password)
	return args.Get(0).(auth.UserCore), args.Error(1)
}

func TestLogin(t *testing.T) {
	mockRepo := new(mockDataAuth)
	svc := service.NewServiceAuth(mockRepo)

	t.Run("success login", func(t *testing.T) {
		expectedUser := auth.UserCore{
			ID:        "123",
			Username:  "john",
			Email:     "john@example.com",
			Password:  "hashed-password",
			Role:      "user",
			Update_At: time.Now(),
		}

		mockRepo.On("Login", "john@example.com", "password123").
			Return(expectedUser, nil).Once()

		result, err := svc.Login("john@example.com", "password123")

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed login", func(t *testing.T) {
		mockRepo.On("Login", "wrong@example.com", "wrongpass").
			Return(auth.UserCore{}, errors.New("invalid credentials")).Once()

		result, err := svc.Login("wrong@example.com", "wrongpass")

		assert.Error(t, err)
		assert.Equal(t, auth.UserCore{}, result)
		mockRepo.AssertExpectations(t)
	})
}
