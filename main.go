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
	// scripts.PopulateDatabase()

	// err := scripts.CreateRandomChallengeMonth(2025, 1) // year=2025, month=1
	// if err != nil {
	// 	fmt.Printf("Failed to create random challenge month: %v\n", err)
	// }

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
	router.Group(func(r chi.Router) {
		r.Use(middle.OptionalAuthMiddleware)

		r.Get("/", controllers.HomeGET)
		r.Get("/leaderboard", controllers.LeaderboardGET)
		r.Get("/day/{year}-{month}-{dayNumber}", controllers.DayGET)
		r.Get("/login", controllers.LoginGET)
		r.Get("/logout", controllers.LogoutGET)
		r.Get("/auth/google", controllers.GoogleGET)
		r.Get("/auth/google/callback", controllers.GoogleCallbackGET)
		r.Get("/auth/apple", controllers.AppleGET)
		r.Post("/auth/apple/callback", controllers.AppleCallbackPOST)
	})

	// Protected routes
	router.Group(func(r chi.Router) {
		r.Use(middle.AuthMiddleware)

		r.Get("/settings", controllers.SettingsGET)
		r.Post("/settings", controllers.SettingsPOST)
		r.Post("/challenge/complete", controllers.CompleteChallengePOST)
	})

	// Start
	http.ListenAndServe(":3000", router)
}
