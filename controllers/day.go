package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/robbyklein/swole/config"
	"github.com/robbyklein/swole/db"
	"github.com/robbyklein/swole/helpers"
	"github.com/robbyklein/swole/sqlc"
)

func DayGET(w http.ResponseWriter, r *http.Request) {
	// Get url parameters
	yearStr := chi.URLParam(r, "year")
	monthStr := chi.URLParam(r, "month")
	dayNumStr := chi.URLParam(r, "dayNumber")

	// Convert them to int32
	yearInt, err := strconv.Atoi(yearStr)
	if err != nil {
		log.Printf("Invalid year: %v\n", err)
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}
	monthInt, err := strconv.Atoi(monthStr)
	if err != nil {
		log.Printf("Invalid month: %v\n", err)
		http.Error(w, "Invalid month", http.StatusBadRequest)
		return
	}
	dayInt, err := strconv.Atoi(dayNumStr)
	if err != nil {
		log.Printf("Invalid day number: %v\n", err)
		http.Error(w, "Invalid day number", http.StatusBadRequest)
		return
	}

	// Get the challenge month
	cm, err := db.Queries.GetChallengeMonthByYearMonth(db.CTX, sqlc.GetChallengeMonthByYearMonthParams{
		Year:  int32(yearInt),
		Month: int32(monthInt),
	})
	if err != nil {
		log.Printf("Challenge month not found for year %d and month %d: %v\n", yearInt, monthInt, err)
		http.Error(w, "Challenge month not found", http.StatusNotFound)
		return
	}

	// Get the day
	dayObj, err := db.Queries.GetDayByMonthIDNumber(db.CTX, sqlc.GetDayByMonthIDNumberParams{
		ChallengeMonthID: cm.ID,
		DayNumber:        int32(dayInt),
	})
	if err != nil {
		log.Printf("Day not found for monthID %d and day %d: %v\n", cm.ID, dayInt, err)
		http.Error(w, "Day not found", http.StatusNotFound)
		return
	}

	// Get the user
	user, loggedIn := r.Context().Value(config.UserContextKey).(sqlc.User)

	// Determine if it's the current day in the user's timezone
	userTZ := helpers.GetUserTimezone(user) // Helper to get user's timezone
	now := time.Now().In(userTZ)
	currentYear, currentMonth, currentDay := now.Date()
	isCurrentDay := currentYear == yearInt && int(currentMonth) == monthInt && currentDay == dayInt

	// List challenges for the day
	dayChallenges, err := db.Queries.ListChallengesForDay(db.CTX, dayObj.ID)
	if err != nil {
		log.Printf("Error listing challenges for dayID %d: %v\n", dayObj.ID, err)
		http.Error(w, "Could not list challenges", http.StatusInternalServerError)
		return
	}

	// Handle empty challenges
	if len(dayChallenges) == 0 {
		log.Printf("No challenges found for dayID %d\n", dayObj.ID)
	}

	// Fetch completed challenges for the user (if logged in)
	completedMap := make(map[int64]bool)
	if loggedIn {
		completedChallenges, err := db.Queries.GetCompletedChallengesForUser(db.CTX, sqlc.GetCompletedChallengesForUserParams{
			UserID: user.ID,
			DayID:  dayObj.ID,
		})
		if err != nil {
			log.Printf("Error fetching completed challenges for userID %d and dayID %d: %v\n", user.ID, dayObj.ID, err)
			http.Error(w, "Could not fetch completed challenges", http.StatusInternalServerError)
			return
		}

		// Populate the map
		for _, challengeID := range completedChallenges {
			completedMap[challengeID] = true
		}
	}

	// Get the month name
	monthName := time.Month(monthInt).String()

	// Prepare data
	data := map[string]interface{}{
		"PageTitle":    "Day",
		"Month":        monthName,
		"DayNumber":    dayInt,
		"DayID":        dayObj.ID,
		"Challenges":   dayChallenges,
		"CompletedMap": completedMap,
		"IsCurrentDay": isCurrentDay,
	}

	// Render the template
	RenderTemplate(w, r, "day", data)
}
