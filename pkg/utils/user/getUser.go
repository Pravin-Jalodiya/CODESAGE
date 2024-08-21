package user

import (
	"cli-project/internal/domain/models"
	"cli-project/pkg/utils/readers"
	"fmt"
)

func GetUser(username string) models.User {
	for i, user := range readers.UserStore {
		if user.Username == username {
			return readers.UserStore[i]
		}
	}

	fmt.Println("user not found")
	return models.User{}
}
