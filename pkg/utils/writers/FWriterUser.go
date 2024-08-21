package writers

import (
	"cli-project/internal/config"
	"cli-project/internal/domain/models"
	"cli-project/pkg/utils/readers"
	"encoding/json"
	"log"
	"os"
)

func FWriterUser(f string, newUser models.User) (bool, error) {

	users := readers.FReaderUser(f, os.O_CREATE|os.O_APPEND|os.O_RDWR)

	users = append(users, newUser)

	// Write updated users back to file
	jsonData, err := json.Marshal(users)
	if err != nil {
		log.Printf("Error marshaling data: %v\n", err)
		return false, err
	}

	err = os.WriteFile(config.USER_FILE, jsonData, 0644)
	if err != nil {
		log.Printf("Error writing to file: %v\n", err)
		return false, err
	}
	readers.UserPassMap[newUser.Username] = newUser.Password
	readers.SyncUserStore()
	return true, nil
}
