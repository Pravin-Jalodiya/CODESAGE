package formatting

import (
	"cli-project/pkg/utils/formatting"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

// TestColorize tests the Colorize function.
func TestColorize(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		fgColor  string
		style    string
		expected string
	}{
		{
			name:     "Red bold text",
			text:     "Hello",
			fgColor:  "red",
			style:    "bold",
			expected: color.New(color.FgRed, color.Bold).Sprint("Hello"),
		},
		{
			name:     "Green underline text",
			text:     "World",
			fgColor:  "green",
			style:    "underline",
			expected: color.New(color.FgGreen, color.Underline).Sprint("World"),
		},
		{
			name:     "Blue plain text",
			text:     "Test",
			fgColor:  "blue",
			style:    "",
			expected: color.New(color.FgBlue).Sprint("Test"),
		},
		{
			name:     "cyan color",
			text:     "cyan",
			fgColor:  "cyan",
			style:    "bold",
			expected: color.New(color.Reset, color.Bold).Sprint("cyan"),
		},
		{
			name:     "Unknown style",
			text:     "Style",
			fgColor:  "yellow",
			style:    "unknown",
			expected: color.New(color.FgYellow).Sprint("Style"),
		},
		{
			name:     "Empty text",
			text:     "",
			fgColor:  "white",
			style:    "",
			expected: color.New(color.FgWhite).Sprint(""),
		},
		{
			name:     "black color",
			text:     "black",
			fgColor:  "black",
			style:    "underline",
			expected: color.New(color.Reset, color.Underline).Sprint("black"),
		},
		{
			name:     "Empty style",
			text:     "EmptyStyle",
			fgColor:  "magenta",
			style:    "",
			expected: color.New(color.FgMagenta).Sprint("EmptyStyle"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatting.Colorize(tt.text, tt.fgColor, tt.style)
			assert.Equal(t, tt.expected, result)
		})
	}
}
