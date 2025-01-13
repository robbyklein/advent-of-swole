package initializers

import (
	"fmt"
	"os"

	"github.com/robbyklein/swole/db"
	"github.com/robbyklein/swole/scripts"
)

func InitDatabase() {
	if os.Getenv("GO_ENV") == "development" {
		// Add challenges
		scripts.AddChallenges()

		// Get challenge months
		challengeMonths, err := db.Queries.ListChallengeMonths(db.CTX)

		if err == nil && len(challengeMonths) == 0 {
			err := scripts.CreateRandomChallengeMonth(2025, 1)

			if err != nil {
				fmt.Printf("Failed to create random challenge month: %v\n", err)
			}
		}

	}
}
