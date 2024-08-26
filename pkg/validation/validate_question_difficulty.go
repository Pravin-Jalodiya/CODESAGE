package validation

import (
	"cli-project/pkg/utils/data_cleaning"
	"errors"
)

func ValidateDifficulty(difficulty string) (string, error) {
	lowerDifficulty := data_cleaning.CleanString(difficulty)
	validDifficulties := map[string]bool{"easy": true, "medium": true, "hard": true}

	if !validDifficulties[lowerDifficulty] {
		return "", errors.New("invalid difficulty level: must be 'easy', 'medium', or 'hard'")
	}
	return lowerDifficulty, nil
}
