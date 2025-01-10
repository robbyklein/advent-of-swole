package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/robbyklein/swole/config"
	"github.com/robbyklein/swole/db"
	"github.com/robbyklein/swole/helpers"
	"github.com/robbyklein/swole/sqlc"
)

func HomeGET(w http.ResponseWriter, r *http.Request) {
	// Fetch the user
	user, _ := r.Context().Value(config.UserContextKey).(sqlc.User)

	// Determine user's timezone
	userTZ := helpers.GetUserTimezone(user)

	// Get the current time in the user's timezone
	currentTime := time.Now().In(userTZ)
	currentYear, currentMonth, currentDay := currentTime.Date()

	fmt.Println(userTZ)
	fmt.Println("CurrentTime")
	fmt.Println(currentTime)

	// Fetch the most recent challenge month
	cm, err := db.Queries.GetMostRecentChallengeMonth(db.CTX)
	if err != nil {
		cm.Year = int32(currentYear)
		cm.Month = int32(currentMonth)
	}

	// Determine the unlocked day
	unlockedDay := calculateUnlockedDay(currentYear, currentMonth, currentDay, cm.Year, cm.Month)

	// Generate days for the homepage
	days := generateDaysSlice(unlockedDay)

	// Calculate "Annual" text
	annual := helpers.OrdinalWords(currentYear - 2024)

	// Calculate time until the next day
	hours, minutes, seconds := calculateTimeUntilNextDay(currentTime, userTZ)

	// Prepare data
	data := map[string]interface{}{
		"PageTitle":           "Homepage",
		"Days":                days,
		"MonthFull":           currentMonth,
		"Month":               cm.Month,
		"Year":                cm.Year,
		"Annual":              annual,
		"HoursUntilNextDay":   hours,
		"MinutesUntilNextDay": minutes,
		"SecondsUntilNextDay": seconds,
	}

	// Render the template
	RenderTemplate(w, r, "home", data)
}

func calculateUnlockedDay(currentYear int, currentMonth time.Month, currentDay int, challengeYear, challengeMonth int32) int32 {
	switch {
	case challengeYear < int32(currentYear),
		(challengeYear == int32(currentYear) && challengeMonth < int32(currentMonth)):
		return 25
	case challengeYear == int32(currentYear) && challengeMonth == int32(currentMonth):
		if currentDay > 25 {
			return 25
		}
		return int32(currentDay)
	default:
		return 25
	}
}

func generateDaysSlice(unlockedDay int32) []map[string]interface{} {
	numbers := helpers.GenerateRangeSlice(1, 25) // [1..25]
	days := make([]map[string]interface{}, len(numbers))

	for i, num := range numbers {
		days[i] = map[string]interface{}{
			"Number":  num,
			"Ordinal": helpers.OrdinalNumbers(num),
			"Active":  int32(num) <= unlockedDay,
		}
	}
	return days
}

func calculateTimeUntilNextDay(currentTime time.Time, userTZ *time.Location) (hours, minutes, seconds int) {
	nextDay := time.Date(
		currentTime.Year(), currentTime.Month(), currentTime.Day()+1,
		0, 0, 0, 0, userTZ,
	)

	duration := nextDay.Sub(currentTime)
	return int(duration.Hours()), int(duration.Minutes()) % 60, int(duration.Seconds()) % 60
}
