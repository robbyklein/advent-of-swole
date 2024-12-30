package helpers

import (
	"context"
	"errors"
	"os"
	"time"

	"google.golang.org/api/idtoken"
)

func VerifyGoogleIDToken(idToken string) (map[string]interface{}, error) {
	// Check for client id
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" {
		return nil, errors.New("GOOGLE_CLIENT_ID is not set")
	}

	// Validate the token
	payload, err := idtoken.Validate(context.Background(), idToken, clientID)
	if err != nil {
		return nil, err
	}

	// Validate expiration
	exp, ok := payload.Claims["exp"].(float64)

	if !ok {
		return nil, errors.New("missing or invalid exp claim in token")
	}

	if time.Now().Unix() > int64(exp) {
		return nil, errors.New("token has expired")
	}

	// Return the valid claims
	return payload.Claims, nil
}
