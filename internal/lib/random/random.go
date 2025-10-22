package random

import (
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewRandomString(lenString int) string {
	b := make([]byte, lenString)
	for i := range b {
		randomIndex := rnd.Intn(len(letters))
		b[i] = letters[randomIndex]
	}
	return string(b)
}
