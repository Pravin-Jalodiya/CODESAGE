package validation

import (
	"testing"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		username string
		expected bool
	}{
		{"username123", true},     // Valid: contains letters and digits after letters
		{"user1name", true},       // Valid: contains letters and digits after letters
		{"123username", false},    // Invalid: starts with digits
		{"124", false},            // Invalid: only digits
		{"112113username", false}, // Invalid: starts with digits
		{"user name", false},      // Valid: contains letters and a space
		{"user name123", false},   // Valid: contains letters, digits, and a space
		{"us3r", true},            // Valid: contains letters and digits after letters
		{"us3r", true},            // Invalid: contains special characters
		{"", false},               // Invalid: empty
		{"user$", false},          // Invalid: contains special characters
		{"user name 123", false},  // Valid: contains letters, digits, and spaces
		{"123user name", false},   // Invalid: starts with digits
	}

	for _, test := range tests {
		result := ValidateUsername(test.username)
		if result != test.expected {
			t.Errorf("ValidateUsername(%q) = %v; want %v", test.username, result, test.expected)
		}
	}
}
