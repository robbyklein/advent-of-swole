package initializers

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

var AppleOauthConfig *oauth2.Config

func CreateAppleOAuthConfig() {
	// Grab env vars
	rootUrl := os.Getenv("ROOT_URL")
	clientId := os.Getenv("APPLE_SERVICES_ID")
	teamId := os.Getenv("APPLE_TEAM_ID")
	keyId := os.Getenv("APPLE_KEY_ID")
	privateKeyPath := os.Getenv("APPLE_PRIVATE_KEY_PATH")

	// Generate the client secret (JWT)
	clientSecret, err := generateAppleClientSecret(teamId, clientId, keyId, privateKeyPath)
	if err != nil {
		panic(fmt.Sprintf("Could not generate Apple client secret: %v", err))
	}

	// Create app oauth client
	AppleOauthConfig = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  rootUrl + "/auth/apple/callback",
		Scopes:       []string{"email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://appleid.apple.com/auth/authorize",
			TokenURL: "https://appleid.apple.com/auth/token",
		},
	}

}

func generateAppleClientSecret(teamId, clientId, keyId, privateKeyPath string) (string, error) {
	// Load the private key file
	keyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read private key file: %v", err)
	}

	// Parse the PEM block
	block, _ := pem.Decode(keyData)
	if block == nil {
		return "", errors.New("failed to parse PEM block containing the private key")
	}

	// Parse the private key (PKCS#8 format)
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse PKCS8 private key: %v", err)
	}

	// Assert that the parsed key is an ECDSA key
	ecdsaKey, ok := parsedKey.(*ecdsa.PrivateKey)
	if !ok {
		return "", errors.New("private key is not an ECDSA key")
	}

	// Generate JWT claims
	now := time.Now()
	claims := jwt.MapClaims{
		"iss": teamId,
		"iat": now.Unix(),
		"exp": now.Add(time.Hour * 24 * 30).Unix(),
		"aud": "https://appleid.apple.com",
		"sub": clientId,
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = keyId // Add Key ID to the header

	// Sign the token with the parsed ECDSA private key
	signedToken, err := token.SignedString(ecdsaKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return signedToken, nil
}
