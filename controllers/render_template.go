package controllers

import (
	"net/http"
	"time"

	"github.com/robbyklein/swole/config"
	"github.com/robbyklein/swole/initializers"
)

func RenderTemplate(w http.ResponseWriter, r *http.Request, templateName string, data map[string]interface{}) {
	// Retrieve the auth session
	authSession, err := initializers.Store.Get(r, config.AUTH_SESSION_KEY)
	if err != nil {
		http.Error(w, "Could not retrieve auth session", http.StatusInternalServerError)
		return
	}

	// Retrieve the flash session
	flashSession, err := initializers.Store.Get(r, config.OTHER_SESSION_KEY)
	if err != nil {
		http.Error(w, "Could not retrieve flash session", http.StatusInternalServerError)
		return
	}

	// Create data if needed
	if data == nil {
		data = make(map[string]interface{})
	}

	// Check for a flash message in the flash session
	if flashMessage, ok := flashSession.Values[config.FLASH_MESSAGE_KEY].(string); ok && flashMessage != "" {
		// Set the flash message in the data map
		data[config.FLASH_MESSAGE_KEY] = flashMessage

		// Clear the flash message
		delete(flashSession.Values, config.FLASH_MESSAGE_KEY)
		if err := flashSession.Save(r, w); err != nil {
			http.Error(w, "Could not save flash session", http.StatusInternalServerError)
			return
		}
	}

	// Check if the user is logged in in the auth session
	if userID, ok := authSession.Values[config.USER_ID_KEY]; ok {
		data["IsLoggedIn"] = true
		data["UserID"] = userID
		if provider, ok := authSession.Values[config.PROVIDER_KEY]; ok {
			data["Provider"] = provider
		}
	} else {
		data["IsLoggedIn"] = false
	}

	// Add the current year for footer
	data["Year"] = time.Now().In(initializers.Location).Year()

	// Render the template
	err = initializers.Templates.ExecuteTemplate(w, templateName, data)
	if err != nil {
		http.Error(w, "Could not render template", http.StatusInternalServerError)
	}
}
