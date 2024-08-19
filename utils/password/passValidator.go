package password

import (
	"unicode"
)

const minPasswordLength = 8 // Minimum allowed length for a strong password

// PasswordValidator evaluates the strength of a password based on several criteria.
// Returns true if the password meets the minimum requirements; otherwise, false.
func ValidatePass(password string) bool {
	if len(password) < minPasswordLength {
		return false
	}

	hasUpper := false
	for _, r := range password {
		if unicode.IsUpper(r) {
			hasUpper = true
			break
		}
	}
	if !hasUpper {
		return false
	}

	hasLower := false
	for _, r := range password {
		if unicode.IsLower(r) {
			hasLower = true
			break
		}
	}
	if !hasLower {
		return false
	}

	hasDigit := false
	for _, r := range password {
		if unicode.IsDigit(r) {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		return false
	}

	hasSpecial := false
	specialChars := []byte("!@#$%^&*()-+?_=,<>/{}[]|`~;")

	// Convert the rune to a byte for comparison
	for _, r := range password {
		if unicode.IsPunct(rune(r)) && specialChars[0] != byte(r) && specialChars[len(specialChars)-1] != byte(r) {
			hasSpecial = true
			break
		}
	}
	if !hasSpecial {
		return false
	}

	return true
}
