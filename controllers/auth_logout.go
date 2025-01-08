package controllers

import (
	"net/http"

	"github.com/robbyklein/swole/config"
	"github.com/robbyklein/swole/helpers"
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

	// Set flash message
	helpers.SetFlashMessage(r, w, "You have been logged out")

	// Send them home
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
