package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/robbyklein/swole/controllers"
	"github.com/robbyklein/swole/db"
	"github.com/robbyklein/swole/helpers"
	"github.com/robbyklein/swole/initializers"
	middle "github.com/robbyklein/swole/middleware"
)

func init() {
	helpers.ValidateEnvVars()
	initializers.LoadIPDatabase()
	initializers.LoadLocation()
	initializers.LoadHTMLTemplates()
	initializers.CreateOAuthConfig()
	initializers.CreateAppleOAuthConfig()
	initializers.InitSessionStore()
	db.ConnectToDatabase()
}

func main() {
	// Setup router
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(httprate.LimitByIP(50, 1*time.Minute))

	// Serve static files
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/", http.FileServer(http.Dir("./public"))).ServeHTTP(w, r)
	})

	// Public routes
	router.Get("/", controllers.HomeGET)
	router.Get("/login", controllers.LoginGET)
	router.Get("/logout", controllers.LogoutGET)
	router.Get("/auth/google", controllers.GoogleGET)
	router.Get("/auth/google/callback", controllers.GoogleCallbackGET)
	router.Get("/auth/apple", controllers.AppleGET)
	router.Get("/auth/apple/callback", controllers.AppleCallbackGET)

	// Protected outes
	router.Group(func(r chi.Router) {
		r.Use(middle.AuthMiddleware)

		r.Get("/settings", controllers.SettingsGET)
		r.Post("/settings", controllers.SettingsPOST)
	})

	// Start
	http.ListenAndServe(":3000", router)
}
