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

	// Generate a random state string
	state := helpers.GenerateRandomString(32)

	// Save the state in the session
	session.Values[config.APPLE_STATE_KEY] = state
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Could not save session", http.StatusInternalServerError)
		return
	}

	// Build the Apple auth URL
	url := initializers.AppleOauthConfig.AuthCodeURL(state)
	url += "&response_mode=form_post"

	// Redirect to auth url
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func AppleCallbackPOST(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Extract state and code from the form
	returnedState := r.FormValue("state")
	code := r.FormValue("code")

	// Validate the code actually exists
	if code == "" {
		http.Error(w, "Missing code parameter", http.StatusBadRequest)
		return
	}

	// Validate the state actually exists
	if returnedState == "" {
		http.Error(w, "Missing state parameter", http.StatusBadRequest)
		return
	}

	// Get the session
	session, err := initializers.Store.Get(r, config.AUTH_SESSION_KEY)
	if err != nil {
		http.Error(w, "Could not get session", http.StatusInternalServerError)
		return
	}

	// Verify states match
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
	userID, userIDExists := claims["sub"].(string)
	if !userIDExists {
		http.Error(w, "User ID not found in ID token", http.StatusInternalServerError)
		return
	}

	// Extract email address
	email, emailExists := claims["email"].(string)
	if !emailExists {
		http.Error(w, "User ID not found in ID token", http.StatusInternalServerError)
		return
	}

	// Generate display name
	displayName := helpers.GenerateDisplayName()

	// Guess timezone
	timezone := helpers.GuessTimezone(r)

	// Create the user
	user, err := db.GetOrCreateUser(ctx, "apple", userID, email, timezone, displayName)
	if err != nil {
		http.Error(w, "User retrieval/creation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Store user information in the session
	session.Values[config.USER_ID_KEY] = user.ID
	session.Values[config.PROVIDER_KEY] = "apple"

	if err := session.Save(r, w); err != nil {
		http.Error(w, "Could not save session", http.StatusInternalServerError)
		return
	}

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
