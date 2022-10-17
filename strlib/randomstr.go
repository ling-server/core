package strlib

import (
	"crypto/rand"

	"github.com/ling-server/core/log"
)

// GenerateRandomStringWithLen generates a random string with length
func GenerateRandomStringWithLen(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	l := len(chars)
	result := make([]byte, length)
	_, err := rand.Read(result)
	if err != nil {
		log.Warningf("Error reading random bytes: %v", err)
	}
	for i := 0; i < length; i++ {
		result[i] = chars[int(result[i])%l]
	}
	return string(result)
}

// GenerateRandomString generate a random string with 32 byte length
func GenerateRandomString() string {
	return GenerateRandomStringWithLen(32)
}
