package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

func GetGravatarURL(email string, size int) string {
	// Trim and normalize the email address
	email = strings.TrimSpace(strings.ToLower(email))

	// Create MD5 hash of the email address
	hash := md5.Sum([]byte(email))
	hashHex := hex.EncodeToString(hash[:])

	// Construct the Gravatar URL
	baseURL := "https://www.gravatar.com/avatar/"
	gravatarURL := fmt.Sprintf("%s%s", baseURL, hashHex)

	// Add size parameter if specified
	if size > 0 {
		gravatarURL = fmt.Sprintf("%s?s=%d", gravatarURL, size)
	}

	return gravatarURL
}
