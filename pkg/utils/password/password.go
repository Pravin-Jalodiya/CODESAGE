package password

import (
	"golang.org/x/crypto/bcrypt"
)

// HashString generates a bcrypt hash for the given password.
func HashString(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// VerifyString verifies if the given password matches the stored hash.
func VerifyString(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
