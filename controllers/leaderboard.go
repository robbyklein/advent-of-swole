package controllers

import (
	"net/http"

	"github.com/robbyklein/swole/db"
)

func LeaderboardGET(w http.ResponseWriter, r *http.Request) {
	// Fetch the leaderboard
	leaderboard, err := db.Queries.GetLeaderboard(db.CTX, 100)
	if err != nil {
		http.Error(w, "Failed to fetch leaderboard: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare data
	data := map[string]interface{}{
		"PageTitle":   "Leaderboard",
		"BodyClass":   "leaderboard",
		"Leaderboard": leaderboard, // Pass the leaderboard data to the template
	}

	// Render the template with data
	RenderTemplate(w, r, "leaderboard", data)
}
