package math

import (
	"math/rand"
	"time"
)

func RandomInt() int {

	rand.NewSource(time.Now().UnixNano())

	randomNumber := rand.Intn(99001) + 1000

	return randomNumber
}
