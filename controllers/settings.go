package controllers

import (
	"net/http"

	"github.com/robbyklein/swole/config"
	"github.com/robbyklein/swole/db"
	"github.com/robbyklein/swole/helpers"
	"github.com/robbyklein/swole/sqlc"
)

func SettingsGET(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(config.UserContextKey).(sqlc.User)

	if !ok {
		http.Error(w, "Must be logged in", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"PageTitle": "Settings",
		"BodyClass": "settings",
		"Timezone":  user.Timezone,
		"Timezones": config.Timezones,
	}

	RenderTemplate(w, r, "settings", data)
}

func SettingsPOST(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(config.UserContextKey).(sqlc.User)

	if !ok {
		http.Error(w, "Must be logged in", http.StatusInternalServerError)
		return
	}

	// Parse the form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Extract and validate form values
	displayName := r.FormValue("displayName")
	timezone := r.FormValue("timezone")

	if displayName == "" || timezone == "" {
		http.Error(w, "Display name and timezone are required", http.StatusBadRequest)
		return
	}

	// Check if the timezone is valid
	isValid := false
	for _, tz := range config.Timezones {
		if tz == timezone {
			isValid = true
			break
		}
	}
	if !isValid {
		http.Error(w, "Invalid timezone selected", http.StatusBadRequest)
		return
	}

	// Update the user
	err := db.Queries.UpdateUser(db.CTX, sqlc.UpdateUserParams{
		ID:          user.ID,
		Timezone:    timezone,
		DisplayName: displayName,
	})

	if err != nil {
		http.Error(w, "Could not update user settings: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve the flash session
	helpers.SetFlashMessage(r, w, "Settings saved successfully!")

	// Redirect back to the settings page with a success message
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
