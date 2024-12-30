package controllers

import (
	"net/http"
)

func HomeGET(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"PageTitle": "Homepage",
	}

	RenderTemplate(w, r, "home", data)
}
