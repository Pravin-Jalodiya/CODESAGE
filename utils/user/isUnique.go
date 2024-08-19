package user

import "cli-project/utils/readers"

func IsUnique(username string) bool {
	if _, exists := readers.UserPassMap[username]; exists {
		return false
	}
	return true
}
