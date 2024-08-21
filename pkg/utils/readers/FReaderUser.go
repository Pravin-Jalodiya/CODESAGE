package readers

import (
	"cli-project/internal/domain/models"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func init() {
	SyncUserStore()
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
