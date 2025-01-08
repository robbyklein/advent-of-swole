package helpers

import (
	"net/http"

	"github.com/robbyklein/swole/config"
	"github.com/robbyklein/swole/db"
	"github.com/robbyklein/swole/sqlc"
)

func GetAuthenticatedUser(r *http.Request) (sqlc.User, bool) {
	// Get id
	userID, ok := GetAuthenticatedUserID(r)
	if !ok {
		return sqlc.User{}, false
	}

	// Get user
	user, err := db.Queries.GetUser(db.CTX, userID)
	if err != nil {
		return sqlc.User{}, false
	}

	return user, true
}

func GetAuthenticatedUserID(r *http.Request) (int64, bool) {
	userID, ok := r.Context().Value(config.UserIDContextKey).(int64)
	return userID, ok
}
