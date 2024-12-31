package db

import (
	"database/sql"
	"errors"

	"github.com/robbyklein/swole/sqlc"
)

func GetOrCreateUser(oauthProvider string, oauthProviderID string) (*sqlc.User, error) {
	// Try to get the user by provider and provider ID
	user, err := Queries.GetUserByProviderId(CTX, sqlc.GetUserByProviderIdParams{
		OauthProvider:   oauthProvider,
		OauthProviderID: oauthProviderID,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Create the user
			user, err := Queries.CreateUser(CTX, sqlc.CreateUserParams{
				OauthProvider:   oauthProvider,
				OauthProviderID: oauthProviderID,
			})

			if err != nil {
				return nil, err
			}

			return &user, nil
		}

		// For any other error, return it
		return nil, err
	}

	// User exists, return it
	return &user, nil
}
