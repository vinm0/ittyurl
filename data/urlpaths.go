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

// Returns the Url instance with the provided source url, if it exists.
// Returns nil, if the path does not exist
func RegisteredPath(source string) (url *Url, found bool) {
	url = UrlBySource(source)

	return url, (url != nil)
}

// Returns a randomized string of legal characters
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

// Returns a random number between min and max, inclusive.
func Rand(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	r := (rand.Int() % (max - min)) + min
	return r
}

// Returns true if path contains none of the prohibited substrings.
func HasProfanity(path []byte) bool {
	for _, word := range StrongFilter {
		if bytes.Contains(path, word) {
			return true
		}
	}
	return false
}
