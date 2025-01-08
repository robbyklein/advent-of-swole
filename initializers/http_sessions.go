package initializers

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

func InitSessionStore() {
	secret := os.Getenv("SESSION_SECRET")

	Store = sessions.NewCookieStore([]byte(secret))

	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24 * 30,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
}
