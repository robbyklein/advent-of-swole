package middleware

import (
	"context"
	"net/http"

	"github.com/robbyklein/swole/config"
	"github.com/robbyklein/swole/initializers"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the auth session
		authSession, err := initializers.Store.Get(r, config.AUTH_SESSION_KEY)
		if err != nil {
			redirectToLoginWithMessage(w, r, "You must be logged in to access that page")
			return
		}

		// Check if the user_id exists in the auth session
		userID, ok := authSession.Values[config.USER_ID_KEY].(int64)
		if !ok || userID == 0 {
			redirectToLoginWithMessage(w, r, "You must be logged in to access that page")
			return
		}

		// Attach userID to the request context
		ctx := context.WithValue(r.Context(), config.UserIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func redirectToLoginWithMessage(w http.ResponseWriter, r *http.Request, message string) {
	// Set a flash message and redirect to the login page
	setFlashMessage(w, r, message)
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

func setFlashMessage(w http.ResponseWriter, r *http.Request, message string) {
	// Retrieve the flash session
	flashSession, err := initializers.Store.Get(r, config.OTHER_SESSION_KEY)
	if err != nil {
		http.Error(w, "Could not retrieve flash session", http.StatusInternalServerError)
		return
	}

	// Set the flash message
	flashSession.Values[config.FLASH_MESSAGE_KEY] = message

	// Save the flash session
	if err := flashSession.Save(r, w); err != nil {
		http.Error(w, "Could not save flash session", http.StatusInternalServerError)
	}
}
