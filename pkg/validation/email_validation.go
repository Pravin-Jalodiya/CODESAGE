package validation

import (
	"regexp"
	"strings"
)

func ValidateEmail(email string) (bool, bool) {

	// Extract the domain from the email
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false, false
	}

	// Define a list of reputable email domains
	reputableDomains := []string{"gmail.com", "outlook.com", "yahoo.com", "watchguard.com", "hotmail.com", "icloud.com"}

	// Regular expression to match a valid email format
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	if !match {
		return false, false
	}

	domain := parts[1]

	// Check if the domain is in the list of reputable domains
	for _, reputableDomain := range reputableDomains {
		if domain == reputableDomain {
			return true, true
		}
	}
	return true, false
}
