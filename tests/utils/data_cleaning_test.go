package utils

import (
	"cli-project/pkg/utils"
	"reflect"
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
			result := utils.CapitalizeWords(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

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
			result := utils.CleanString(tt.input)
			if result != tt.expected {
				t.Errorf("CleanString(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCleanTags tests the CleanTags function.
func TestCleanTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Normal tags",
			input:    "  Tag1, Tag2,Tag3  ",
			expected: []string{"tag1", "tag2", "tag3"},
		},
		{
			name:     "Tags with extra spaces",
			input:    "  TagA  ,  TagB , TagC  ",
			expected: []string{"taga", "tagb", "tagc"},
		},
		{
			name:     "Single tag",
			input:    " TagOnly ",
			expected: []string{"tagonly"},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: []string{""},
		},
		{
			name:     "Only commas",
			input:    ",,,",
			expected: []string{"", "", "", ""},
		},
		{
			name:     "Tags with special characters",
			input:    " Tag1# , @Tag2,Tag3! ",
			expected: []string{"tag1#", "@tag2", "tag3!"},
		},
		{
			name:     "String with no commas",
			input:    "NoCommaTag",
			expected: []string{"nocommatag"},
		},
		{
			name:     "Multiple consecutive commas",
			input:    "Tag1,,,Tag2,,Tag3",
			expected: []string{"tag1", "", "", "tag2", "", "tag3"},
		},
		{
			name:     "Tags with numbers",
			input:    "Tag1,Tag2,1234",
			expected: []string{"tag1", "tag2", "1234"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.CleanTags(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("CleanTags(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
