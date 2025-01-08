package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/robbyklein/swole/sqlc"
)

func GetOrCreateChallenge(ctx context.Context, challenge sqlc.CreateChallengeParams) (*sqlc.Challenge, error) {
	existingChallenge, err := Queries.GetChallengeByDescription(ctx, challenge.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newChallenge, createErr := Queries.CreateChallenge(ctx, challenge)
			if createErr != nil {
				return nil, createErr
			}
			return &newChallenge, nil
		}

		return nil, err
	}

	return &existingChallenge, nil
}

func LinkChallengeToDay(ctx context.Context, dayID int64, challengeID int64) error {
	return Queries.LinkChallengeToDay(ctx, sqlc.LinkChallengeToDayParams{
		DayID:       dayID,
		ChallengeID: challengeID,
	})
}
