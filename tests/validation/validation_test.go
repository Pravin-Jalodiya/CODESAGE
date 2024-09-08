package validation

import (
	"cli-project/pkg/validation"
	"strings"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	cases := []struct {
		email       string
		isValid     bool
		isReputable bool
		description string
	}{
		{"test@gmail.com", true, true, "Valid reputable email"},
		{"test@unknown.com", true, false, "Valid non-reputable email"},
		{"invalidemail@", false, false, "Email without domain"},
		{"@gmail.com", false, false, "Email without local part"},
		{"invalidemail", false, false, "Email without @ symbol"},
		{"", false, false, "Empty email"},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			isValid, isReputable := validation.ValidateEmail(c.email)
			if isValid != c.isValid || isReputable != c.isReputable {
				t.Errorf("Expected (%v, %v) for email %s, but got (%v, %v)",
					c.isValid, c.isReputable, c.email, isValid, isReputable)
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
	cases := []struct {
		name    string
		isValid bool
	}{
		{"John", true},
		{"A", false},
		{"", false},
		{"Valid Name", true},
		{"Invalid@Name", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := validation.ValidateName(c.name)
			if result != c.isValid {
				t.Errorf("Expected %v, but got %v", c.isValid, result)
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
		{name: "password lenght less than minimum allowed", input: "vsmall", expected: false},
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
		{"invalid character", "%$$$%@$#$@--", false},
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
			result, _ := validation.ValidateQuestionDifficulty(tt.input)
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
		{"Valid Leetcode link", "https://Leetcode.com/problems/example-problem/", "https://leetcode.com/problems/example-problem/"},
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

func TestValidateTitleSlug(t *testing.T) {
	tests := []struct {
		titleSlug   string
		expectValid bool
		expectError bool
	}{
		{"valid-title-slug", true, false},   // Valid title slug
		{"", false, true},                   // Empty title slug
		{"another-valid-slug", true, false}, // Another valid slug
		{" ", false, true},                  // Space only slug
	}

	for _, test := range tests {
		valid, err := validation.ValidateTitleSlug(test.titleSlug)
		if valid != test.expectValid {
			t.Errorf("Expected validity for slug '%s' to be %v, got %v", test.titleSlug, test.expectValid, valid)
		}
		if (err != nil) != test.expectError {
			t.Errorf("Expected error for slug '%s' to be %v, got %v", test.titleSlug, test.expectError, err != nil)
		}
	}
}
