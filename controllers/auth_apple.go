package controllers

import (
	"net/http"

	"github.com/robbyklein/swole/config"
	"github.com/robbyklein/swole/db"
	"github.com/robbyklein/swole/helpers"
	"github.com/robbyklein/swole/initializers"
)

func AppleGET(w http.ResponseWriter, r *http.Request) {
	// Create or retrieve the session
	session, err := initializers.Store.Get(r, config.AUTH_SESSION_KEY)
	if err != nil {
		http.Error(w, "Could not get session", http.StatusInternalServerError)
		return
	}

	// Generate a random state
	state := helpers.GenerateRandomString(32)

	// Save the state in the session
	session.Values[config.APPLE_STATE_KEY] = state
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Could not save session", http.StatusInternalServerError)
		return
	}

	// Build the Apple auth URL
	url := initializers.AppleOauthConfig.AuthCodeURL(state)

	// Redirect to Apple
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func AppleCallbackGET(w http.ResponseWriter, r *http.Request) {
	// Extract state and code from url
	returnedState := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	// Validate the code actually exists
	if code == "" {
		http.Error(w, "Missing code parameter", http.StatusBadRequest)
		return
	}

	// Get the session
	session, err := initializers.Store.Get(r, config.AUTH_SESSION_KEY)
	if err != nil {
		http.Error(w, "Could not get session", http.StatusInternalServerError)
		return
	}

	// Verify state matches
	savedState, ok := session.Values[config.APPLE_STATE_KEY].(string)
	if !ok || savedState == "" || savedState != returnedState {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	// Exchange the code for a token
	ctx := r.Context()
	token, err := initializers.AppleOauthConfig.Exchange(ctx, code)
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
	claims, err := helpers.VerifyAppleIDToken(idToken)
	if err != nil {
		http.Error(w, "Could not verify ID token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract user identifier
	userID := claims["sub"].(string)
	provider := "apple"

	// Retrieve or create the user in the database
	user, err := db.GetOrCreateUser(provider, userID)
	if err != nil {
		http.Error(w, "Could not create or retrieve user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a settings row for the user
	displayName := helpers.GenerateDisplayName()
	timezone := helpers.GuessTimezone(r)

	_, err = db.CreateSettings(user.ID, timezone, displayName)
	if err != nil {
		http.Error(w, "Could not create settings for user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Store user information in the session
	session.Values[config.USER_ID_KEY] = user.ID
	session.Values[config.PROVIDER_KEY] = "apple"

	if err := session.Save(r, w); err != nil {
		http.Error(w, "Could not save session", http.StatusInternalServerError)
		return
	}

	// Redirect to home
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
