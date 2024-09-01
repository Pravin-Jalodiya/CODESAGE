package validation

import (
	"cli-project/pkg/validation"
	"strings"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid email with reputable domain", "example@gmail.com", true},
		{"Valid email with reputable domain", "example@yahoo.com", true},
		{"Invalid email format", "invalid_email", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := validation.ValidateEmail(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateEmail(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateCountryName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid country", "united states", true},
		{"Invalid country", "Unknown Country", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := validation.ValidateCountryName(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateCountryName(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid name", "Pravin k", true},
		{"Empty name", "", false},
		{"Name too short", "A", false},
		{"Name too long", strings.Repeat("a", 46), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validation.ValidateName(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateName(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateOrganizationName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid organization name", "Example Corporation", true},
		{"Invalid organization name (too short)", "X", false},
		{"Invalid organization name (too long)", strings.Repeat("a", 46), false},
		{"Invalid organization name (contains special characters)", "Example!Corp", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := validation.ValidateOrganizationName(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateOrganizationName(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid password", "StrongP@ssword123!", true},
		{"Weak password (no uppercase)", "weakpassword123!", false},
		{"Weak password (no lowercase)", "WEAKPASSWORD123!", false},
		{"Weak password (no digit)", "StrongP@ssword", false},
		{"Weak password (no special character)", "StrongPassword123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validation.ValidatePassword(tt.input)
			if result != tt.expected {
				t.Errorf("ValidatePassword(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid username", "JohnDoe123", true},
		{"Invalid username (too short)", "JD", false},
		{"Invalid username (too long)", strings.Repeat("a", 46), false},
		{"Invalid username (no letter)", "123456789", false},
		{"Valid username (no digit after letter)", "John", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validation.ValidateUsername(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateUsername(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateDifficulty(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Valid difficulty", "easy", "easy"},
		{"Valid difficulty", "medium", "medium"},
		{"Valid difficulty", "hard", "hard"},
		{"Invalid difficulty", "invalid", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := validation.ValidateDifficulty(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateDifficulty(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateQuestionLink(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Valid Leetcode link", "https://Leetcode.com/problems/example-problem/", "https://Leetcode.com/problems/example-problem/"},
		{"Invalid URL", "invalid-url", ""},
		{"Leetcode link without scheme", "Leetcode.com/problems/example-problem/", ""},
		{"Leetcode link without host", "https://example.com/problems/example-problem/", ""},
		{"Non-Leetcode link", "https://example.com", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := validation.ValidateQuestionLink(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateQuestionLink(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateQuestionID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid positive integer ID", "10", true},
		{"Invalid negative ID", "-1", false},
		{"Invalid non-integer ID", "abc", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := validation.ValidateQuestionID(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateQuestionID(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
