package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/robbyklein/swole/sqlc"
)

func GetOrCreateUser(ctx context.Context, oauthProvider, oauthProviderID, email, timezone, displayName string) (*sqlc.User, error) {
	// Attempt to retrieve the user by provider and provider ID
	user, err := Queries.GetUserByProviderId(ctx, sqlc.GetUserByProviderIdParams{
		OauthProvider:   oauthProvider,
		OauthProviderID: oauthProviderID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			user, err = Queries.CreateUser(ctx, sqlc.CreateUserParams{
				OauthProvider:   oauthProvider,
				OauthProviderID: oauthProviderID,
				Email:           email,
				Timezone:        timezone,
				DisplayName:     displayName,
			})
			if err != nil {
				return nil, err
			}
			return &user, nil
		}

		return nil, err
	}

	return &user, nil
}
