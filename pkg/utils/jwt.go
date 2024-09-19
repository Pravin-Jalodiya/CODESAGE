package utils

import (
	"cli-project/internal/config"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// CreateJwtToken generates a JWT token for the given user ID and role
func CreateJwtToken(username string, userId string, role string) (string, error) {
	// Define JWT claims
	claims := jwt.MapClaims{
		"username": username,
		"userId":   userId,
		"role":     role,
		"exp":      time.Now().Add(time.Minute).Unix(), // Token expiry time (1 minute)
	}

	// Create a new JWT token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(config.SECRET_KEY)
	if err != nil {
		// Log the error if needed (uncomment the line below)
		// logger.Logger.Errorw("Error signing token", "error", err, "time", time.Now())
		return "", errors.New("error creating jwt token")
	}

	return tokenString, nil
}
