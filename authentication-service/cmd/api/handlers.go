package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJson(w, r, &requestPayload) // Read the payload into requestPayload struct
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	// Validate that the user exists in the db data
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJson(w, errors.New("Invalid credentials"), http.StatusBadRequest)
		return
	}

	validPassword, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !validPassword {
		app.errorJson(w, errors.New("Invalid credentials"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Loggedin user %s", user.Email), // In real app here we would return the token as well
		Data:    user,
	}

	app.writeJson(w, http.StatusAccepted, payload)
}
