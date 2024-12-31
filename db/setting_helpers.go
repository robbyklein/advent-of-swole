package db

import (
	"database/sql"
	"errors"

	"github.com/robbyklein/swole/sqlc"
)

func CreateSettings(userID int64, timezone string, displayName string) (*sqlc.Setting, error) {
	// Check if settings already exist for the user
	settings, err := Queries.GetSettings(CTX, userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Create new settings for the user
			settings, err := Queries.CreateSettings(CTX, sqlc.CreateSettingsParams{
				UserID:      userID,
				Timezone:    timezone,
				DisplayName: displayName,
			})

			if err != nil {
				return nil, err
			}

			return &settings, nil
		}

		// For any other error, return it
		return nil, err
	}

	// Settings already exist, return them
	return &settings, nil
}
