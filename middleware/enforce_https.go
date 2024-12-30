package middleware

import "net/http"

func EnforceHTTPS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Forwarded-Proto") != "https" {
			http.Error(w, "HTTPS required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
