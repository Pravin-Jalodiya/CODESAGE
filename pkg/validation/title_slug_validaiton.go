package validation

import (
	"fmt"
	"strings"
)

func ValidateTitleSlug(titleSlug string) (bool, error) {
	// Add validation logic for titleSlug, such as format checks
	titleSlug = strings.TrimSpace(titleSlug)
	if len(titleSlug) == 0 {
		return false, fmt.Errorf("title slug cannot be empty")
	}
	// Add more validations as needed
	return true, nil
}