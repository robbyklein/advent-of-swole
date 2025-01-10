package controllers

import (
	"net/http"

	"github.com/robbyklein/swole/db"
	"github.com/robbyklein/swole/helpers"
)

type LeaderboardEntry struct {
	UserID      int64  `json:"user_id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	TotalPoints int32  `json:"total_points"`
	GravatarURL string `json:"gravatar_url"`
}

func LeaderboardGET(w http.ResponseWriter, r *http.Request) {
	// Fetch the leaderboard data
	leaderboardRows, err := db.Queries.GetLeaderboard(db.CTX, 100)
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
			GravatarURL: helpers.GetGravatarURL(row.Email, 320), // Compute Gravatar URL
		}
	}

	// Prepare data for rendering
	data := map[string]interface{}{
		"PageTitle":   "Leaderboard",
		"Leaderboard": leaderboard, // Use the enriched leaderboard
	}

	// Render the template
	RenderTemplate(w, r, "leaderboard", data)
}
