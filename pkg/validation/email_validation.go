package validation

import (
	"regexp"
)

// ValidateEmail checks if the email format is valid
func ValidateEmail(email string) bool {
	// Simple email regex
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	return match
}
