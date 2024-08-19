package readers

import (
	"cli-project/models"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var (
	UserStore   []models.User
	UserPassMap = map[string]string{}
)

func init() {
	SyncUser()
	for _, user := range UserStore {
		UserPassMap[user.Username] = user.Password
	}
}

func FReaderUser(f string, flag int) []models.User {

	var users []models.User

	file, err := os.OpenFile(f, flag, 0644)
	if err != nil {
		fmt.Println("Error opening file")
		//log error concurrently
	}

	byteValue, _ := io.ReadAll(file)
	err = json.Unmarshal(byteValue, &users)
	if err != nil {
		return nil
	}

	return users
}
