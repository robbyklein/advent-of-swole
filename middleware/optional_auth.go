package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/robbyklein/swole/config"
	"github.com/robbyklein/swole/db"
	"github.com/robbyklein/swole/initializers"
)

func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the auth session
		authSession, err := initializers.Store.Get(r, config.AUTH_SESSION_KEY)
		if err == nil {
			// Get the user id
			if userID, ok := authSession.Values[config.USER_ID_KEY].(int64); ok && userID != 0 {
				// Get the user
				user, err := db.Queries.GetUser(db.CTX, userID)
				if err == nil {
					// Attach to context
					ctx := context.WithValue(r.Context(), config.UserContextKey, user)
					r = r.WithContext(ctx)
				} else {
					fmt.Printf("Error retrieving user: %v\n", err)
				}
			}
		} else {
			fmt.Printf("Error retrieving session: %v\n", err)
		}
		next.ServeHTTP(w, r)
	})
}
