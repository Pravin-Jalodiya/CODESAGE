package data_cleaning

import (
	"cli-project/pkg/utils/data_cleaning"
	"testing"
)

// TestCleanString tests the CleanString function.
func TestCleanString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal string",
			input:    "  Hello World  ",
			expected: "hello world",
		},
		{
			name:     "All uppercase",
			input:    "  GOLANG  ",
			expected: "golang",
		},
		{
			name:     "Mixed case with spaces",
			input:    "  GoLang iS aWesoMe  ",
			expected: "golang is awesome",
		},
		{
			name:     "Already clean",
			input:    "clean",
			expected: "clean",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "String with only spaces",
			input:    "     ",
			expected: "",
		},
		{
			name:     "String with spaces around",
			input:    "   Test   ",
			expected: "test",
		},
		{
			name:     "String with special characters",
			input:    "   @#$%^&*()_+   ",
			expected: "@#$%^&*()_+",
		},
		{
			name:     "String with numbers and letters",
			input:    "  1234 ABCD efgh  ",
			expected: "1234 abcd efgh",
		},
		{
			name:     "Multiline string",
			input:    "  This is \nA Test  ",
			expected: "this is \na test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := data_cleaning.CleanString(tt.input)
			if result != tt.expected {
				t.Errorf("CleanString(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}
