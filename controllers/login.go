package controllers

import (
	"net/http"
)

func LoginGET(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"PageTitle": "Login",
		"BodyClass": "login",
	}

	RenderTemplate(w, r, "login", data)

}
