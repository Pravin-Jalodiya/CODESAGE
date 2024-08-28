package data_cleaning

import (
	"strings"
	"unicode"
)

// CapitalizeWords capitalizes the first letter of each word in a string.
func CapitalizeWords(s string) string {
	words := strings.Fields(s) // Split the string into words
	for i, word := range words {
		// Capitalize the first letter of each word
		if len(word) > 0 {
			words[i] = string(unicode.ToUpper(rune(word[0]))) + word[1:]
		}
	}
	return strings.Join(words, " ")
}
