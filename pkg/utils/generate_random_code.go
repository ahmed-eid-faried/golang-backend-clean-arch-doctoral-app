package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomCode() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(9000) + 1000 // Ensures a 4-digit number between 1000 and 9999
}
