package utils

import (
	"cli-project/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// TestHashString tests the HashString function.
func TestHashString(t *testing.T) {
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
			hash, err := utils.HashString(tt.password)
			if err != nil {
				t.Fatalf("HashString() returned an error: %v", err)
			}

			// Ensure the hash is not empty
			if hash == "" {
				t.Errorf("HashString() returned an empty hash")
			}

			// Check if the hash is valid
			err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(tt.password))
			if err != nil {
				t.Errorf("HashString() returned an invalid hash: %v", err)
			}
		})
	}
}

// TestVerifyString tests the VerifyString function.
func TestVerifyString(t *testing.T) {
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
			result := utils.VerifyString(tt.password, tt.hash)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function to generate a hash for a password
func getHashForPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hash)
}
