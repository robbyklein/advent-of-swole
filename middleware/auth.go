package middleware

import (
	"net/http"

	"github.com/robbyklein/swole/constants"
	"github.com/robbyklein/swole/initializers"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the session
		session, err := initializers.Store.Get(r, constants.SESSION_KEY)

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		// Check if the user_id exists in the session
		userID, ok := session.Values[constants.USER_ID_KEY]

		if !ok || userID == nil {
			session.Values[constants.FLASH_MESSAGE_KEY] = "You must be logged in to access this page"

			if err := session.Save(r, w); err != nil {
				http.Error(w, "Could not save session", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		// User is authenticated, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
