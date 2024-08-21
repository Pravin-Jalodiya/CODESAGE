package validation

import "cli-project/pkg/utils/readers"

func IsUsernameUnique(username string) bool {
	if _, exists := readers.UserPassMap[username]; exists {
		return false
	}
	return true
}
