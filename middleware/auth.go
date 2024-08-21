package middleware

import "cli-project/pkg/globals"

func Auth(userID string) {
	globals.ActiveUser = userID
}
