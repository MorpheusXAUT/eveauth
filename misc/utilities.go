package misc

import (
	"math/rand"
)

func GenerateRandomString(length int) string {
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, length)

	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}
