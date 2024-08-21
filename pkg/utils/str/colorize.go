package str

import (
	"github.com/fatih/color"
)

// Define foreground and background colors
var (
	colors = map[string]color.Attribute{
		"black":   color.FgBlack,
		"red":     color.FgRed,
		"green":   color.FgGreen,
		"yellow":  color.FgYellow,
		"blue":    color.FgBlue,
		"magenta": color.FgMagenta,
		"cyan":    color.FgCyan,
		"white":   color.FgWhite,
	}
	styles = map[string]color.Attribute{
		"bold":      color.Bold,
		"underline": color.Underline,
	}
)

// Colorize generates a colorized str with the specified foreground, background, and style
func Colorize(text, fg, bg, style string) string {
	fgColor, fgOk := colors[fg]
	bgColor, bgOk := colors[bg]
	styleColor, styleOk := styles[style]

	if !fgOk {
		fgColor = color.Reset
	}
	if !bgOk {
		bgColor = color.Reset
	}
	if !styleOk {
		styleColor = color.Reset
	}

	return color.New(fgColor, bgColor, styleColor).SprintFunc()(text)
}
