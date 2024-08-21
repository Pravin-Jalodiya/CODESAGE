package validation

import (
	"strings"
	"unicode"
)

// ValidateName checks if the provided name is valid
func ValidateName(name string) bool {
	// Trim any extra whitespace from the start and end of the name
	name = strings.TrimSpace(name)

	// Check for maximum length
	if len(name) > 45 {
		return false
	}

	// Check if name contains only letters and spaces
	for _, r := range name {
		if !(unicode.IsLetter(r) || unicode.IsSpace(r)) {
			return false
		}
	}

	return true
}
