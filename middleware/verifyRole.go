package middleware

import (
	usr "cli-project/utils/user"
	"slices"
)

func VerifyRole(username string, roles ...string) bool {

	user := usr.GetUser(username)

	return slices.Contains(roles, user.Role)

}
