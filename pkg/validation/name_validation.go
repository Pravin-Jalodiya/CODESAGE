package validation

// ValidateName checks if the name is valid
func ValidateName(name string) bool {
	// Name must be non-empty and only contain alphabetic characters and spaces
	if len(name) <= 2 || len(name) >= 31 {
		return false
	}

	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == ' ') {
			return false
		}
	}
	return true
}
