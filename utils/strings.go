package utils

import (
	"crypto/rand"
	"io"
	"math/big"
)

type CharType string

const (
	SPECIAL         CharType = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ^$?!@#%&.,-_"
	ALPHA_NUMERIC   CharType = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NUMERIC         CharType = "0123456789"
	LETTERS         CharType = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CAPITAL_LETTERS CharType = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	SMALL_LETTERS   CharType = "abcdefghijklmnopqrstuvwxyz"
)

func GenerateString(length int, characters CharType) string {
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = characters[int(b[i])%len(characters)]
	}
	return string(b)
}

func GenerateRandomString(length int) string {
	return GenerateString(length, SPECIAL)
}

func GenerateStringRandomLength(min int, max int) string {
	val, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	length := min + int(val.Int64())
	return GenerateString(length, SPECIAL)
}
