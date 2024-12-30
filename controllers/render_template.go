package controllers

import (
	"net/http"

	"github.com/robbyklein/swole/constants"
	"github.com/robbyklein/swole/initializers"
)

func RenderTemplate(w http.ResponseWriter, r *http.Request, templateName string, data map[string]interface{}) {
	// Add flash message from session
	session, err := initializers.Store.Get(r, constants.SESSION_KEY)

	if err != nil {
		http.Error(w, "Could not retrieve session", http.StatusInternalServerError)
		return
	}

	// Create data if needed
	if data == nil {
		data = make(map[string]interface{})
	}

	// Check for a flash message
	if flashMessage, ok := session.Values[constants.FLASH_MESSAGE_KEY].(string); ok && flashMessage != "" {

		// Set it
		data[constants.FLASH_MESSAGE_KEY] = flashMessage

		// Clear it
		delete(session.Values, constants.FLASH_MESSAGE_KEY)

		if err := session.Save(r, w); err != nil {
			http.Error(w, "Could not save session", http.StatusInternalServerError)
			return
		}
	}

	// Check if user is logged in
	if userID, ok := session.Values[constants.USER_ID_KEY]; ok {
		data["IsLoggedIn"] = true
		data["UserID"] = userID
		if provider, ok := session.Values[constants.PROVIDER_KEY]; ok {
			data["Provider"] = provider
		}
	} else {
		data["IsLoggedIn"] = false
	}

	// Render the template
	err = initializers.Templates.ExecuteTemplate(w, templateName, data)

	if err != nil {
		http.Error(w, "Could not render template", http.StatusInternalServerError)
	}
}
