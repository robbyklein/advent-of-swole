package helpers

import (
	"log"
	"os"
)

// ValidateEnvVars ensures all required environment variables are set
func ValidateEnvVars() {
	requiredVars := []string{
		"ROOT_URL",
		"DATABASE_URL",
		"TEST_DATABASE_URL",
		"GOOGLE_CLIENT_ID",
		"GOOGLE_CLIENT_SECRET",
		"SESSION_SECRET",
		"APPLE_SERVICES_ID",
		"APPLE_TEAM_ID",
		"APPLE_KEY_ID",
		"APPLE_PRIVATE_KEY_PATH",
	}

	for _, envVar := range requiredVars {
		if value := os.Getenv(envVar); value == "" {
			log.Fatalf("Missing required environment variable: %s", envVar)
		}
	}
}
