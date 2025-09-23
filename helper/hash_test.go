package helper_test

import (
	"go_rest_native_sekolah/helper"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestCheckPassword(t *testing.T) {
	t.Run("Password valid", func(t *testing.T) {
		// Arrange
		password := "mypassword"
		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		// Act
		result := helper.CheckPassword(password, string(hash))

		// Assert
		assert.True(t, result, "Password harus cocok")
	})

	t.Run("Password tidak valid", func(t *testing.T) {
		// Arrange
		password := "mypassword"
		wrongPassword := "wrongpass"
		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		// Act
		result := helper.CheckPassword(wrongPassword, string(hash))

		// Assert
		assert.False(t, result, "Password tidak boleh cocok")
	})
}
