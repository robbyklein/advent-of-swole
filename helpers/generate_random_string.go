package helpers

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomString(length int) string {
	b := make([]byte, length)

	if _, err := rand.Read(b); err != nil {
		return "fallback-state"
	}

	return base64.URLEncoding.EncodeToString(b)
}
