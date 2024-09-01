package data_cleaning

import (
	"cli-project/pkg/utils/data_cleaning"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapitalizeWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single word lowercase",
			input:    "hello",
			expected: "Hello",
		},
		{
			name:     "multiple words lowercase",
			input:    "hello world",
			expected: "Hello World",
		},
		{
			name:     "mixed case input",
			input:    "hElLo wOrLd",
			expected: "HElLo WOrLd",
		},
		{
			name:     "all caps input",
			input:    "HELLO WORLD",
			expected: "HELLO WORLD",
		},
		{
			name:     "leading spaces",
			input:    "  hello world",
			expected: "Hello World",
		},
		{
			name:     "trailing spaces",
			input:    "hello world  ",
			expected: "Hello World",
		},
		{
			name:     "leading and trailing spaces",
			input:    "  hello world  ",
			expected: "Hello World",
		},
		{
			name:     "multiple spaces between words",
			input:    "hello   world",
			expected: "Hello World",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single letter words",
			input:    "a b c d e",
			expected: "A B C D E",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := data_cleaning.CapitalizeWords(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
