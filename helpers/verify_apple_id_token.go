package helpers

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AppleJWKResponse struct {
	Keys []AppleJWK `json:"keys"`
}

type AppleJWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func getAppleKeys() ([]AppleJWK, error) {
	resp, err := http.Get("https://appleid.apple.com/auth/keys")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jwkResp AppleJWKResponse
	if err := json.Unmarshal(body, &jwkResp); err != nil {
		return nil, err
	}

	return jwkResp.Keys, nil
}

func appleKeyToPublicKey(jwk AppleJWK) (*rsa.PublicKey, error) {
	if jwk.Kty != "RSA" {
		return nil, fmt.Errorf("unsupported key type: %s", jwk.Kty)
	}

	// Decode base64URL modulus
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode n (modulus): %w", err)
	}
	modulus := new(big.Int).SetBytes(nBytes)

	// Decode base64URL exponent
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode e (exponent): %w", err)
	}

	// Convert exponent bytes to int
	var expBytes []byte

	if len(eBytes) < 8 {
		pad := make([]byte, 8-len(eBytes))
		expBytes = append(pad, eBytes...)
	} else {
		expBytes = eBytes
	}
	exponent := binary.BigEndian.Uint64(expBytes)

	pubKey := &rsa.PublicKey{
		N: modulus,
		E: int(exponent),
	}
	return pubKey, nil
}

func VerifyAppleIDToken(idToken string) (map[string]interface{}, error) {
	clientID := os.Getenv("APPLE_SERVICES_ID")

	if clientID == "" {
		return nil, errors.New("APPLE_SERVICES_ID not set")
	}

	parsedToken, err := jwt.Parse(idToken, func(t *jwt.Token) (interface{}, error) {
		// 1. Extract "kid" from the token header
		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, errors.New("no kid found in token header")
		}

		// 2. Fetch Apple’s JWKS
		appleKeys, err := getAppleKeys()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch apple keys: %w", err)
		}

		// 3. Find the matching key
		var appleKey AppleJWK
		found := false
		for _, k := range appleKeys {
			if k.Kid == kid {
				appleKey = k
				found = true
				break
			}
		}
		if !found {
			return nil, errors.New("no matching Apple public key found for kid")
		}

		// 4. Convert JWK -> rsa.PublicKey
		pubKey, err := appleKeyToPublicKey(appleKey)
		if err != nil {
			return nil, fmt.Errorf("failed to convert Apple JWK to public key: %w", err)
		}

		// 5. Validate the token’s alg
		if t.Method.Alg() != appleKey.Alg {
			return nil, fmt.Errorf("unexpected token signing method: %v", t.Header["alg"])
		}

		// Return the correct public key so the library can verify the signature.
		return pubKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to verify token signature: %w", err)
	}

	// Check the overall validity (signature check, etc.)
	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract the claims as a MapClaims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to parse token claims")
	}

	// Check standard claims (issuer, audience, expiry, etc.)
	if iss, ok := claims["iss"].(string); !ok || iss != "https://appleid.apple.com" {
		return nil, errors.New("invalid issuer")
	}

	if aud, ok := claims["aud"].(string); !ok || aud != clientID {
		return nil, fmt.Errorf("invalid audience, expected %s got %s", clientID, aud)
	}

	// Validate expiration
	if exp, ok := claims["exp"].(float64); !ok || time.Now().Unix() > int64(exp) {
		return nil, errors.New("token has expired")
	}

	// If everything checks out, return the claims.
	return claims, nil
}
