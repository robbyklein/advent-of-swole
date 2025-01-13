package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/robbyklein/swole/db"
	"github.com/robbyklein/swole/helpers"
	"github.com/robbyklein/swole/sqlc"
)

type LeaderboardEntry struct {
	UserID      int64  `json:"user_id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	TotalPoints int32  `json:"total_points"`
	GravatarURL string `json:"gravatar_url"`
}

func LeaderboardGET(w http.ResponseWriter, r *http.Request) {
	// Extract year and month from URL parameters
	year := chi.URLParam(r, "year")
	month := chi.URLParam(r, "month")

	// Validate year and month
	yearInt, err := strconv.Atoi(year)
	if err != nil || yearInt < 2000 || yearInt > 2100 {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	monthInt, err := strconv.Atoi(month)
	if err != nil || monthInt < 1 || monthInt > 12 {
		http.Error(w, "Invalid month", http.StatusBadRequest)
		return
	}

	// Fetch the challenge_month_id based on the year and month
	challengeMonth, err := db.Queries.GetChallengeMonthByYearMonth(db.CTX, sqlc.GetChallengeMonthByYearMonthParams{
		Year:  int32(yearInt),
		Month: int32(monthInt),
	})
	if err != nil {
		http.Error(w, "Challenge month not found: "+err.Error(), http.StatusNotFound)
		return
	}

	// Fetch leaderboard data for the specified challenge month
	leaderboardRows, err := db.Queries.GetLeaderboard(db.CTX, sqlc.GetLeaderboardParams{
		ChallengeMonthID: challengeMonth.ID,
		Limit:            100,
	})
	if err != nil {
		http.Error(w, "Failed to fetch leaderboard: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Map to custom struct and compute Gravatar URLs
	leaderboard := make([]LeaderboardEntry, len(leaderboardRows))
	for i, row := range leaderboardRows {
		leaderboard[i] = LeaderboardEntry{
			UserID:      row.UserID,
			DisplayName: row.DisplayName,
			Email:       row.Email,
			TotalPoints: row.TotalPoints,
			GravatarURL: helpers.GetGravatarURL(row.Email, 320),
		}
	}

	// Prepare data for rendering
	data := map[string]interface{}{
		"PageTitle":   fmt.Sprintf("Leaderboard - %04d-%02d", yearInt, monthInt),
		"Leaderboard": leaderboard,
	}

	// Render the template
	RenderTemplate(w, r, "leaderboard", data)
}
