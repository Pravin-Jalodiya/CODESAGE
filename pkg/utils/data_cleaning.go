package utils

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

func CleanString(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func CleanTags(tags string) []string {
	tagList := strings.Split(tags, ",")
	for i, tag := range tagList {
		tagList[i] = CleanString(tag)
	}
	return tagList
}
