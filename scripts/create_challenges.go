package scripts

import (
	"log"

	"github.com/robbyklein/swole/challenges"
	"github.com/robbyklein/swole/db"
	"github.com/robbyklein/swole/sqlc"
)

func AddChallenges() {
	allChallenges := challenges.GetAllChallenges()

	for _, challenge := range allChallenges {
		params := sqlc.CreateChallengeParams{
			Description:            challenge.Description,
			DescriptionMetric:      challenge.DescriptionMetric,
			Category:               challenge.Category,
			MuscleGroups:           challenge.MuscleGroups,
			Difficulty:             int32(challenge.Difficulty),
			CaloriesBurnedEstimate: int32(challenge.CaloriesBurnedEstimate),
		}

		db.GetOrCreateChallenge(db.CTX, params)

	}

	log.Println("Database populated with challenges.")
}
