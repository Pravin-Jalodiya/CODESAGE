package middleware

import (
	usr "cli-project/pkg/utils/user"
	"slices"
)

func VerifyRole(username string, roles ...string) bool {

	user := usr.GetUser(username)

	return slices.Contains(roles, user.Role)

}
