package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (c *Config) AuthHandler(w http.ResponseWriter, r *http.Request) {
	// set a json response
	var ReqPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	log.Println("Received request in auth service",
		r.Method,
		r.URL,
		r.RemoteAddr,
	)

	// store json to data
	err := c.readJSON(w, r, &ReqPayload)
	if err != nil {
		log.Println("err in auth service, while storing json as data", err)
		c.errorJSON(w, err, http.StatusBadRequest)

	}

	// ------ authenticate the user ------
	// get the user by email from DB
	user, err := c.Models.User.GetByEmail(ReqPayload.Email)
	if err != nil {
		c.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}
	// if user exists in DB, validate the password
	valid, err := user.PasswordMatches(ReqPayload.Password)
	if err != nil || !valid {
		c.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}
	// if user is authenticated with password, send ok status
	payload := JsonResponse{
		ErrorPresent: false,
		Data:         user,
		Message:      fmt.Sprintf("User %s logged in successfully", user.Email),
	}

	c.writeJSON(w, http.StatusAccepted, payload)
}
