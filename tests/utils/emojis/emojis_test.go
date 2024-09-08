package emojis

import (
	"cli-project/pkg/utils/emojis"
	"testing"
)

func TestEmojis(t *testing.T) {
	cases := []struct {
		name  string
		emoji string
	}{
		{"Signup", emojis.Signup},
		{"Login", emojis.Login},
		{"Exit", emojis.Exit},
		{"Error", emojis.Error},
		{"Success", emojis.Success},
		{"Profile", emojis.Profile},
		{"Stats", emojis.Stats},
		{"Settings", emojis.Settings},
		{"Question", emojis.Question},
		{"Info", emojis.Info},
		{"Back", emojis.Back},
		{"View", emojis.View},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.emoji == "" {
				t.Error("Emoji should not be empty")
			}
		})
	}
}
