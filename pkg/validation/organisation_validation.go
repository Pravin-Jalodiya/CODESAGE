package validation

import (
	"errors"
	"regexp"
)

func ValidateOrganizationName(orgName string) (bool, error) {

	if len(orgName) <= 1 || len(orgName) > 40 {
		return false, errors.New("invalid organization name : name must be between 2 and 40 characters")
	}
	// Regex to allow only letters and spaces
	const orgNameRegex = `^[a-zA-Z\s]+$`
	match, _ := regexp.MatchString(orgNameRegex, orgName)

	if !match {
		return false, errors.New("invalid organization name : only letters and spaces are allowed")
	}

	return true, nil
}
