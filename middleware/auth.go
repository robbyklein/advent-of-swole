package middleware

import (
	"net/http"

	"github.com/robbyklein/swole/config"
	"github.com/robbyklein/swole/initializers"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the auth session
		authSession, err := initializers.Store.Get(r, config.AUTH_SESSION_KEY)
		if err != nil {
			// Redirect to login if the auth session is not retrievable
			setFlashMessage(w, r, "You must be logged in to access that page")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		// Check if the user_id exists in the auth session
		userID, ok := authSession.Values[config.USER_ID_KEY]
		if !ok || userID == nil {
			// Set a flash message in the flash session
			setFlashMessage(w, r, "You must be logged in to access that page")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		// User is authenticated, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

// Helper function to set a flash message in the flash session
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
