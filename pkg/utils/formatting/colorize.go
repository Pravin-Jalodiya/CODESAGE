package formatting

import (
	"github.com/fatih/color"
)

// Colorize generates a colorized string with the specified foreground color and style
func Colorize(text, fgColor, style string) string {
	var c *color.Color

	// Apply foreground color
	switch fgColor {
	case "black":
		c = color.New(color.FgBlack)
	case "red":
		c = color.New(color.FgRed)
	case "green":
		c = color.New(color.FgGreen)
	case "yellow":
		c = color.New(color.FgYellow)
	case "blue":
		c = color.New(color.FgBlue)
	case "magenta":
		c = color.New(color.FgMagenta)
	case "cyan":
		c = color.New(color.FgCyan)
	case "white":
		c = color.New(color.FgWhite)
	default:
		c = color.New(color.Reset)
	}

	// Apply style
	switch style {
	case "bold":
		c = c.Add(color.Bold)
	case "underline":
		c = c.Add(color.Underline)
	}

	// Return the formatted string
	return c.Sprint(text)
}
