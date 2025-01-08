package controllers

import (
	"net/http"
	"time"

	"github.com/robbyklein/swole/config"
	"github.com/robbyklein/swole/helpers"
	"github.com/robbyklein/swole/initializers"
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

	// Retrieve userID from the context
	if userID, ok := helpers.GetAuthenticatedUserID(r); ok {
		// User is logged in, add user-related data
		data["IsLoggedIn"] = true
		data["UserID"] = userID
	} else {
		// User is not logged in
		data["IsLoggedIn"] = false
	}

	// Add the current year for the footer
	data["Year"] = time.Now().In(initializers.Location).Year()

	// Render the template
	err = initializers.Templates.ExecuteTemplate(w, templateName, data)
	if err != nil {
		http.Error(w, "Could not render template", http.StatusInternalServerError)
	}
}
