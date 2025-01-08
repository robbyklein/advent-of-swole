package helpers

import (
	"net/http"

	"github.com/robbyklein/swole/config"
	"github.com/robbyklein/swole/initializers"
)

func SetFlashMessage(r *http.Request, w http.ResponseWriter, msg string) {
	// Retrieve the flash session
	flashSession, err := initializers.Store.Get(r, config.OTHER_SESSION_KEY)
	if err != nil {
		http.Error(w, "Could not get flash session", http.StatusInternalServerError)
		return
	}

	// Add the flash message
	flashSession.Values[config.FLASH_MESSAGE_KEY] = msg
	if err := flashSession.Save(r, w); err != nil {
		http.Error(w, "Could not save flash message", http.StatusInternalServerError)
		return
	}
}
