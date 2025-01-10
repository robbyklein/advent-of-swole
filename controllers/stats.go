package controllers

import (
	"net/http"

	"github.com/robbyklein/swole/config"
	"github.com/robbyklein/swole/db"
	"github.com/robbyklein/swole/sqlc"
)

func StatsGET(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(config.UserContextKey).(sqlc.User)

	if !ok {
		http.Error(w, "Must be logged in", http.StatusInternalServerError)
		return
	}

	calories, err := db.Queries.GetCaloriesStats(db.CTX, sqlc.GetCaloriesStatsParams{
		ChallengeMonthID: 1,
		UserID:           user.ID,
	})

	if err == nil {
		http.Error(w, "Failed to get statistics", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"PageTitle": "Stats",
	}

	RenderTemplate(w, r, "stats", data)
}
