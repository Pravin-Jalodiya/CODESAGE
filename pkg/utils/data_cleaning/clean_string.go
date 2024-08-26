package data_cleaning

import "strings"

func CleanString(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}
