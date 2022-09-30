package util

import (
	"crypto/rand"
	"math/big"
)

func RandomString(n int) string {
	letterRunes := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	result := make([]rune, n)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))
		result[i] = letterRunes[num.Int64()]
	}

	return string(result)
}
