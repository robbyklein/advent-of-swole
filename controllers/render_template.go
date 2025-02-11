package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/robbyklein/swole/config"
	"github.com/robbyklein/swole/helpers"
	"github.com/robbyklein/swole/initializers"
	"github.com/robbyklein/swole/sqlc"
)

func RenderTemplate(w http.ResponseWriter, r *http.Request, templateName string, data map[string]interface{}) {
	// Retrieve the flash session
	flashSession, err := initializers.Store.Get(r, config.OTHER_SESSION_KEY)
	if err != nil {
		http.Error(w, "Could not retrieve flash session", http.StatusInternalServerError)
		return
	}

	// Initialize the data map if it's nil
	if data == nil {
		data = make(map[string]interface{})
	}

	// Add flash message to the data map if it exists
	if flashMessage, ok := flashSession.Values[config.FLASH_MESSAGE_KEY].(string); ok && flashMessage != "" {
		data[config.FLASH_MESSAGE_KEY] = flashMessage

		// Clear the flash message
		delete(flashSession.Values, config.FLASH_MESSAGE_KEY)
		if err := flashSession.Save(r, w); err != nil {
			http.Error(w, "Could not save flash session", http.StatusInternalServerError)
			return
		}
	}

	// Retrieve user from the context
	user, ok := r.Context().Value(config.UserContextKey).(sqlc.User)
	if ok {
		gravatar := helpers.GetGravatarURL(user.Email, 320)

		data["Gravatar"] = gravatar
		data["IsLoggedIn"] = true
		data["User"] = user
	} else {
		data["IsLoggedIn"] = false
	}

	// Add the current year for the footer
	data["FooterYear"] = time.Now().In(initializers.Location).Year()
	data["CurrentChallenge"] = "2025-01"

	// Render the template
	err = initializers.Templates.ExecuteTemplate(w, templateName, data)
	if err != nil {
		fmt.Printf("Error rendering template: %v\n", err)
		http.Error(w, "Could not render template", http.StatusInternalServerError)
		return
	}
}
