package validation

import (
	"unicode"
)

// ValidateUsername checks if the username is valid
func ValidateUsername(username string) bool {
	if len(username) == 0 {
		return false
	}

	hasLetter := false
	hasDigitAfterLetter := false

	for _, r := range username {
		if unicode.IsLetter(r) {
			hasLetter = true
			// Digits after a letter are allowed
			hasDigitAfterLetter = true
		} else if unicode.IsDigit(r) {
			if hasLetter {
				hasDigitAfterLetter = true
			} else {
				// Digits are not allowed if there has been no letter before
				return false
			}
		} else {
			// Invalid character found
			return false
		}
	}

	return hasLetter && hasDigitAfterLetter
}
