package helpers

import (
	"time"

	"github.com/robbyklein/swole/sqlc"
)

func GetUserTimezone(user sqlc.User) *time.Location {
	// Define the default timezone
	const defaultTimezone = "America/Los_Angeles"

	// Load the default location
	defaultLocation, err := time.LoadLocation(defaultTimezone)
	if err != nil {
		return time.UTC
	}

	// Load and return the user's specified timezone
	if user.ID != 0 {
		userLocation, err := time.LoadLocation(user.Timezone)
		if err != nil {
			return defaultLocation
		}

		return userLocation
	}

	return defaultLocation

}
