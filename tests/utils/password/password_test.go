package password

import (
	"cli-project/pkg/utils/password"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// TestHashPassword tests the HashPassword function.
func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "Normal password",
			password: "mysecurepassword",
		},
		{
			name:     "Empty password",
			password: "",
		},
		{
			name:     "Short password",
			password: "short",
		},
		{
			name:     "Special characters",
			password: "!@#$%^&*()_+",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := password.HashPassword(tt.password)
			if err != nil {
				t.Fatalf("HashPassword() returned an error: %v", err)
			}

			// Ensure the hash is not empty
			if hash == "" {
				t.Errorf("HashPassword() returned an empty hash")
			}

			// Check if the hash is valid
			err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(tt.password))
			if err != nil {
				t.Errorf("HashPassword() returned an invalid hash: %v", err)
			}
		})
	}
}

// TestVerifyPassword tests the VerifyPassword function.
func TestVerifyPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		hash     string
		expected bool
	}{
		{
			name:     "Correct password",
			password: "mysecurepassword",
			hash:     getHashForPassword("mysecurepassword"),
			expected: true,
		},
		{
			name:     "Incorrect password",
			password: "wrongpassword",
			hash:     getHashForPassword("mysecurepassword"),
			expected: false,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     getHashForPassword("mysecurepassword"),
			expected: false,
		},
		{
			name:     "Empty hash",
			password: "mysecurepassword",
			hash:     "",
			expected: false,
		},
		{
			name:     "Empty password and hash",
			password: "",
			hash:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := password.VerifyPassword(tt.password, tt.hash)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function to generate a hash for a password
func getHashForPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash)
}
