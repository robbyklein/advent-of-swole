package middleware

import (
	"context"
	"net/http"

	"github.com/robbyklein/swole/config"
	"github.com/robbyklein/swole/initializers"
)

func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the auth session
		authSession, err := initializers.Store.Get(r, config.AUTH_SESSION_KEY)
		if err == nil {
			// Check if the user_id exists in the auth session
			if userID, ok := authSession.Values[config.USER_ID_KEY].(int64); ok && userID != 0 {
				// Attach userID to the request context
				ctx := context.WithValue(r.Context(), config.UserIDContextKey, userID)
				r = r.WithContext(ctx)
			}
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
