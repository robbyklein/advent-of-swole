package controllers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/robbyklein/swole/config"
	"github.com/robbyklein/swole/db"
	"github.com/robbyklein/swole/helpers"
	"github.com/robbyklein/swole/sqlc"
)

func CompleteChallengePOST(w http.ResponseWriter, r *http.Request) {
	// Get the user
	user, ok := r.Context().Value(config.UserContextKey).(sqlc.User)
	if !ok {
		http.Error(w, "Must be logged in", http.StatusInternalServerError)
		return
	}

	// Parse challenge_id from form
	challengeID, err := strconv.ParseInt(r.FormValue("challenge_id"), 10, 64)
	if err != nil || challengeID <= 0 {
		http.Error(w, "Invalid challenge_id", http.StatusBadRequest)
		return
	}

	// Parse day_id from form
	dayID, err := strconv.ParseInt(r.FormValue("day_id"), 10, 64)
	if err != nil || dayID <= 0 {
		http.Error(w, "Invalid day_id", http.StatusBadRequest)
		return
	}

	// Fetch the day from the database
	dayObj, err := db.Queries.GetDay(db.CTX, dayID)
	if err != nil {
		http.Error(w, "Day not found", http.StatusBadRequest)
		return
	}

	// Fetch the challenge month for the day
	challengeMonth, err := db.Queries.GetChallengeMonth(db.CTX, dayObj.ChallengeMonthID)
	if err != nil {
		http.Error(w, "Challenge month not found", http.StatusInternalServerError)
		return
	}

	// Determine if we are in development mode
	isDev := os.Getenv("GO_ENV") == "development"

	// Only enforce the "current day" check if NOT in development mode
	if !isDev {
		userTZ := helpers.GetUserTimezone(user)
		now := time.Now().In(userTZ)
		currentYear, currentMonth, currentDay := now.Date()

		if !(int(currentYear) == int(challengeMonth.Year) &&
			int(currentMonth) == int(challengeMonth.Month) &&
			int(currentDay) == int(dayObj.DayNumber)) {
			http.Error(w, "You can only complete challenges for the current day in your timezone", http.StatusForbidden)
			return
		}
	}

	// Record the challenge completion in the database
	params := sqlc.CompleteChallengeParams{
		UserID:      user.ID,
		ChallengeID: challengeID,
		DayID:       dayID,
	}

	err = db.Queries.CompleteChallenge(db.CTX, params)
	if err != nil {
		http.Error(w, "Failed to record challenge completion: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect the user back to the referring page (or home if none)
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}
