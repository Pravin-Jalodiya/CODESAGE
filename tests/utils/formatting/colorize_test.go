package formatting

import (
	"cli-project/pkg/utils/formatting"
	"testing"
)

// TestColorize tests the Colorize function.

func TestColorize(t *testing.T) {
	cases := []struct {
		text    string
		fgColor string
		style   string
	}{
		{"text", "black", "bold"},
		{"text", "red", "underline"},
		{"text", "green", ""},
		{"text", "yellow", "bold"},
		{"text", "blue", "underline"},
		{"text", "magenta", ""},
		{"text", "cyan", "bold"},
		{"text", "white", "underline"},
		{"text", "invalidcolor", "invalidstyle"},
	}

	for _, c := range cases {
		t.Run(c.fgColor+"_"+c.style, func(t *testing.T) {
			result := formatting.Colorize(c.text, c.fgColor, c.style)
			if result == "" {
				t.Error("Formatted text should not be empty")
			}
		})
	}
}
