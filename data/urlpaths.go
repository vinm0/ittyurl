package data

import (
	"bytes"
	"math/rand"
	"time"
)

const (
	LEGAL_CHARACTERS = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890-"

	MIN_PATH_CHAR = 3
	MAX_PATH_CHAR = 10
)

func RandomPath() string {
	pathLen := Rand(MIN_PATH_CHAR, MAX_PATH_CHAR)
	path := make([]byte, pathLen)

	for i := 0; i < pathLen; i++ {
		ch := Rand(0, len(LEGAL_CHARACTERS))
		path = append(path, LEGAL_CHARACTERS[ch])
	}

	for HasProfanity(path) {
		RandomPath()
	}

	return string(path)
}

func Rand(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	r := (rand.Int() % (max - min)) + min
	return r
}

func HasProfanity(path []byte) bool {
	for _, word := range StrongFilter {
		if bytes.Contains(path, word) {
			return true
		}
	}
	return false
}
