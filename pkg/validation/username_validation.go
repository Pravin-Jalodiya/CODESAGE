package validation

// ValidateName checks if the name is valid
func ValidateUsername(name string) bool {
	// Name must be non-empty and only contain alphabetic characters and spaces
	if len(name) == 0 {
		return false
	}
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == ' ') {
			return false
		}
	}
	return true
}
