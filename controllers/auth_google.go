package controllers

import (
	"net/http"

	"github.com/robbyklein/swole/constants"
	"github.com/robbyklein/swole/helpers"
	"github.com/robbyklein/swole/initializers"
	"golang.org/x/oauth2"
)

func GoogleGET(w http.ResponseWriter, r *http.Request) {
	// Create or retrieve the session
	session, err := initializers.Store.Get(r, constants.SESSION_KEY)

	if err != nil {
		http.Error(w, "Could not get session", http.StatusInternalServerError)
		return
	}

	// Generate a random state string
	state := helpers.GenerateRandomString(32)

	// Save it in the session
	session.Values[constants.GOOGLE_STATE_KEY] = state

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
	session, err := initializers.Store.Get(r, constants.SESSION_KEY)

	if err != nil {
		http.Error(w, "Could not get session", http.StatusInternalServerError)
		return
	}

	// Extract code and state from url
	code := r.URL.Query().Get("code")
	returnedState := r.URL.Query().Get("state")

	// Get the saved session state
	savedState, ok := session.Values[constants.GOOGLE_STATE_KEY].(string)

	// Verify state matches
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

	// Extract user identifier
	userID := claims["sub"].(string)

	// Store user information in the session
	session.Values[constants.USER_ID_KEY] = userID
	session.Values[constants.PROVIDER_KEY] = "google"

	if err := session.Save(r, w); err != nil {
		http.Error(w, "Could not save session", http.StatusInternalServerError)
		return
	}

	// Redirect to the home page
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
