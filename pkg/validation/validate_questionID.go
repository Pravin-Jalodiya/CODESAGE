package validation

import (
	"errors"
	"strconv"
)

func ValidateQuestionID(questionID string) (string, error) {
	qid, err := strconv.Atoi(questionID)
	if err != nil || qid <= 0 {
		return "", errors.New("invalid question ID must be a positive number")
	}
	return questionID, nil
}
