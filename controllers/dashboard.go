package controllers

import (
	"net/http"
)

func DashboardGET(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"PageTitle": "Dashboard",
	}

	RenderTemplate(w, r, "dashboard", data)
}
