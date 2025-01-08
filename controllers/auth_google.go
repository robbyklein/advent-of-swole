package controllers

import (
	"net/http"

	"github.com/robbyklein/swole/config"
	"github.com/robbyklein/swole/db"
	"github.com/robbyklein/swole/helpers"
	"github.com/robbyklein/swole/initializers"
	"golang.org/x/oauth2"
)

func GoogleGET(w http.ResponseWriter, r *http.Request) {
	// Create or retrieve the session
	session, err := initializers.Store.Get(r, config.AUTH_SESSION_KEY)
	if err != nil {
		http.Error(w, "Could not get session", http.StatusInternalServerError)
		return
	}

	// Generate a random state string
	state := helpers.GenerateRandomString(32)

	// Save it in the session
	session.Values[config.GOOGLE_STATE_KEY] = state
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Could not save session", http.StatusInternalServerError)
		return
	}

	// Build the Google auth URL
	url := initializers.GoogleOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)

	// Redirect to it
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallbackGET(w http.ResponseWriter, r *http.Request) {
	// Get the session
	session, err := initializers.Store.Get(r, config.AUTH_SESSION_KEY)
	if err != nil {
		http.Error(w, "Could not get session", http.StatusInternalServerError)
		return
	}

	// Extract code and state from url
	code := r.URL.Query().Get("code")
	returnedState := r.URL.Query().Get("state")

	// Verify state matches
	savedState, ok := session.Values[config.GOOGLE_STATE_KEY].(string)
	if !ok || savedState == "" || savedState != returnedState {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	// Exchange the code for a token
	ctx := r.Context()
	token, err := initializers.GoogleOauthConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Could not exchange code for token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the ID token
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "Could not extract ID token", http.StatusInternalServerError)
		return
	}

	// Decode and validate the ID token to extract user info
	claims, err := helpers.VerifyGoogleIDToken(idToken)
	if err != nil {
		http.Error(w, "Could not verify ID token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Grab the provider id
	userID, userIDExists := claims["sub"].(string)
	if !userIDExists {
		http.Error(w, "User ID not found in ID token", http.StatusInternalServerError)
		return
	}

	// Grab their email address
	email, emailExists := claims["email"].(string)
	if !emailExists {
		http.Error(w, "Email address not found in ID token", http.StatusInternalServerError)
		return
	}

	// Generate a display name
	displayName := helpers.GenerateDisplayName()

	// Guess the users timezone
	timezone := helpers.GuessTimezone(r)

	// Retrieve or create the user
	user, err := db.GetOrCreateUser(ctx, "google", userID, email, timezone, displayName)
	if err != nil {
		http.Error(w, "User retrieval/creation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Store user information in the session
	session.Values[config.USER_ID_KEY] = user.ID
	session.Values[config.PROVIDER_KEY] = "google"
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Could not save session", http.StatusInternalServerError)
		return
	}

	// Redirect to the home page
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
