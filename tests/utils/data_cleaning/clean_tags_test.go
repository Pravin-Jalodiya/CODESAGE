package data_cleaning

import (
	"cli-project/pkg/utils/data_cleaning"
	"reflect"
	"testing"
)

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
			result := data_cleaning.CleanTags(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("CleanTags(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
