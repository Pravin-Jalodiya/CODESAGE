package utils

import (
	"cli-project/pkg/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestConvertToIST tests the ConvertToIST function.
func TestConvertToIST(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "Start of year",
			input:    time.Date(2024, time.January, 1, 0, 0, 0, 0, time.Local),
			expected: "01/01/2024 00:00:00",
		},
		{
			name:     "End of year",
			input:    time.Date(2024, time.December, 31, 23, 59, 59, 0, time.Local),
			expected: "31/12/2024 23:59:59",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ConvertToIST(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
