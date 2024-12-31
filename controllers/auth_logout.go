package controllers

import (
	"net/http"

	"github.com/robbyklein/swole/config"
	"github.com/robbyklein/swole/initializers"
)

func LogoutGET(w http.ResponseWriter, r *http.Request) {
	// Retrieve the auth session
	authSession, err := initializers.Store.Get(r, config.AUTH_SESSION_KEY)
	if err != nil {
		http.Error(w, "Could not get auth session", http.StatusInternalServerError)
		return
	}

	// Invalidate the auth session
	authSession.Options.MaxAge = -1
	if err := authSession.Save(r, w); err != nil {
		http.Error(w, "Could not invalidate auth session", http.StatusInternalServerError)
		return
	}

	// Retrieve the flash session
	flashSession, err := initializers.Store.Get(r, config.OTHER_SESSION_KEY)
	if err != nil {
		http.Error(w, "Could not get flash session", http.StatusInternalServerError)
		return
	}

	// Add the flash message
	flashSession.Values[config.FLASH_MESSAGE_KEY] = "You have been logged out"
	if err := flashSession.Save(r, w); err != nil {
		http.Error(w, "Could not save flash message", http.StatusInternalServerError)
		return
	}

	// Redirect to the home or login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
