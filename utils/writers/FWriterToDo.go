package writers

import (
	"cli-project/config"
	"cli-project/models"
	"cli-project/utils/readers"
	"encoding/json"
	"log"
	"os"
)

func FWriterToDo(f string, newUser []models.User) (bool, error) {

	users := readers.FReaderUser(f, os.O_CREATE|os.O_APPEND|os.O_RDWR)

	users = newUser

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
	readers.SyncUser()
	return true, nil
}
