package utils

import (
	"cli-project/pkg/utils"
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
			result := utils.Colorize(c.text, c.fgColor, c.style)
			if result == "" {
				t.Error("Formatted text should not be empty")
			}
		})
	}
}

func TestEmojis(t *testing.T) {
	cases := []struct {
		name  string
		emoji string
	}{
		{"Signup", utils.SignupEmoji},
		{"Login", utils.LoginEmoji},
		{"Exit", utils.ExitEmoji},
		{"Error", utils.ErrorEmoji},
		{"Success", utils.SuccessEmoji},
		{"Profile", utils.ProfileEmoji},
		{"Stats", utils.StatsEmoji},
		{"Settings", utils.SettingsEmoji},
		{"Question", utils.QuestionEmoji},
		{"Info", utils.InfoEmoji},
		{"Back", utils.BackEmoji},
		{"View", utils.ViewEmoji},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.emoji == "" {
				t.Error("Emoji should not be empty")
			}
		})
	}
}
