package scripts

import (
	"log"

	"github.com/robbyklein/swole/challenges"
	"github.com/robbyklein/swole/db"
	"github.com/robbyklein/swole/sqlc"
)

// PopulateDatabase is an example function that:
// 1) Connects to the DB (via db.ConnectToDatabase())
// 2) Retrieves challenges from your 'challenges' package
// 3) Inserts or updates them (just as an example)
func PopulateDatabase() {
	// 1) Connect to DB
	db.ConnectToDatabase()
	ctx := db.CTX

	// 2) Get all challenges
	allChallenges := challenges.GetAllChallenges()

	// 3) Insert or update them in the DB
	for _, challenge := range allChallenges {
		params := sqlc.CreateChallengeParams{
			Description:            challenge.Description,
			Category:               challenge.Category,
			MuscleGroups:           challenge.MuscleGroups,
			Difficulty:             int32(challenge.Difficulty),
			CaloriesBurnedEstimate: int32(challenge.CaloriesBurnedEstimate),
		}

		db.GetOrCreateChallenge(ctx, params)

	}

	log.Println("PopulateDatabase completed!")
}
