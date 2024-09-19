package utils

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

func ReadCSV(filePath string) ([][]string, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("error opening question file")
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("error closing file")
		}
	}(file)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}
