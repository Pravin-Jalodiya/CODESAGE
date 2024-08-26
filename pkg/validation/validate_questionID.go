package validation

import (
	"errors"
	"strconv"
)

func ValidateQuestionID(questionID string) (bool, error) {
	qid, err := strconv.Atoi(questionID)
	if err != nil || qid <= 0 {
		return false, errors.New("invalid question ID : must be a positive number")
	}
	return true, nil
}
