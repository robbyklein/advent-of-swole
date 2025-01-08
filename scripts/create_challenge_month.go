package scripts

import (
	"log"
	"math/rand"
	"time"

	"github.com/robbyklein/swole/db"
	"github.com/robbyklein/swole/sqlc"
)

// CreateRandomChallengeMonth:
//  1. Connects to the DB (if not already).
//  2. Creates a challenge_month for the given year and month.
//  3. Creates 25 days (day_number = 1..25) in that challenge_month.
//  4. For each day, assigns 1..3 random challenges from the existing "challenges" table.
func CreateRandomChallengeMonth(year, month int32) error {
	// 1) Connect to DB (if not already connected)
	db.ConnectToDatabase()
	ctx := db.CTX
	queries := db.Queries

	// 2) Create the challenge_month
	cm, err := queries.CreateChallengeMonth(ctx, sqlc.CreateChallengeMonthParams{
		Month: month,
		Year:  year,
	})
	if err != nil {
		return err
	}
	log.Printf("Created challenge_month ID=%d (year=%d, month=%d)\n", cm.ID, cm.Year, cm.Month)

	// 3) Create 25 days
	dayIDs := make([]int64, 0, 25)
	for d := int32(1); d <= 25; d++ {
		day, err := queries.CreateDay(ctx, sqlc.CreateDayParams{
			ChallengeMonthID: cm.ID,
			DayNumber:        d,
		})
		if err != nil {
			return err
		}
		dayIDs = append(dayIDs, day.ID)
		log.Printf("Created day %d (ID=%d) for challenge_month %d\n", d, day.ID, cm.ID)
	}

	// 4) Load all challenges so we can pick random ones
	allChallenges, err := queries.ListChallenges(ctx)
	if err != nil {
		return err
	}
	if len(allChallenges) == 0 {
		log.Println("No challenges found in DB; skipping random assignment.")
		return nil
	}

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// For each day, pick 1..3 random challenges and link
	for _, dayID := range dayIDs {
		// pick a random number from 1..3
		numChallenges := rand.Intn(3) + 1 // 1, 2, or 3

		// Optionally: track used indexes to avoid duplicates.
		// For simplicity, we won't here. It's okay if a day gets the same challenge more than once
		// (unless your schema prevents duplicates via a primary key constraint).

		for i := 0; i < numChallenges; i++ {
			// pick random challenge index
			idx := rand.Intn(len(allChallenges))
			challenge := allChallenges[idx]

			// Link them
			err := queries.LinkChallengeToDay(ctx, sqlc.LinkChallengeToDayParams{
				DayID:       dayID,
				ChallengeID: challenge.ID,
			})
			if err != nil {
				log.Printf(
					"Error linking challenge ID=%d to day ID=%d: %v\n",
					challenge.ID,
					dayID,
					err,
				)
			} else {
				log.Printf("Linked challenge ID=%d to day ID=%d\n", challenge.ID, dayID)
			}
		}
	}

	log.Println("CreateRandomChallengeMonth completed successfully!")
	return nil
}
