package controllers

import (
	"net/http"

	"github.com/robbyklein/swole/config"
	"github.com/robbyklein/swole/db"
	"github.com/robbyklein/swole/initializers"
	"github.com/robbyklein/swole/sqlc"
)

func SettingsGET(w http.ResponseWriter, r *http.Request) {
	// Retrieve the auth session
	authSession, err := initializers.Store.Get(r, config.AUTH_SESSION_KEY)
	if err != nil {
		http.Error(w, "Could not retrieve auth session", http.StatusInternalServerError)
		return
	}

	// Get the user ID from the session
	userID, ok := authSession.Values[config.USER_ID_KEY].(int64)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// Fetch the user's settings from the database
	settings, err := db.Queries.GetSettings(r.Context(), userID)
	if err != nil {
		http.Error(w, "Could not retrieve user settings: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Pass settings data to the template
	data := map[string]interface{}{
		"PageTitle":   "Settings",
		"BodyClass":   "settings",
		"DisplayName": settings.DisplayName,
		"Timezone":    settings.Timezone,
		"IsLoggedIn":  true,
		"UserID":      userID,
		"Provider":    authSession.Values[config.PROVIDER_KEY],
		"Timezones":   config.Timezones,
	}

	RenderTemplate(w, r, "settings", data)
}

func SettingsPOST(w http.ResponseWriter, r *http.Request) {
	// Retrieve the auth session
	authSession, err := initializers.Store.Get(r, config.AUTH_SESSION_KEY)
	if err != nil {
		http.Error(w, "Could not retrieve auth session", http.StatusInternalServerError)
		return
	}

	// Get the user ID from the session
	userID, ok := authSession.Values[config.USER_ID_KEY].(int64)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
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

	// Update the user's settings in the database
	err = db.Queries.UpdateSettings(r.Context(), sqlc.UpdateSettingsParams{
		UserID:      userID,
		DisplayName: displayName,
		Timezone:    timezone,
	})
	if err != nil {
		http.Error(w, "Could not update user settings: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve the flash session
	flashSession, err := initializers.Store.Get(r, config.OTHER_SESSION_KEY)
	if err != nil {
		http.Error(w, "Could not get flash session", http.StatusInternalServerError)
		return
	}

	// Add the flash message
	flashSession.Values[config.FLASH_MESSAGE_KEY] = "Settings saved successfully!"
	if err := flashSession.Save(r, w); err != nil {
		http.Error(w, "Could not save flash message", http.StatusInternalServerError)
		return
	}

	// Redirect back to the settings page with a success message
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
