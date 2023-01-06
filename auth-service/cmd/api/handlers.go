package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type JSONPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (c *Config) AuthHandler(w http.ResponseWriter, r *http.Request) {
	// set a json response
	var reqPayload JSONPayload

	log.Println("Received request in auth service",
		r.Method,
		r.URL,
		r.RemoteAddr,
	)

	// store json to data
	err := c.readJSON(w, r, &reqPayload)
	if err != nil {
		log.Println("err in auth service, while storing json as data", err)
		c.errorJSON(w, err, http.StatusBadRequest)

	}

	// ------ authenticate the user ------
	// get the user by email from DB
	user, err := c.Models.User.GetByEmail(reqPayload.Email)
	if err != nil {
		c.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}
	// if user exists in DB, validate the password
	valid, err := user.PasswordMatches(reqPayload.Password)
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

	// log user authenticated event
	l := LogPayload{
		Name: "Event - User Log In",
		Data: user.Email + " logged in successfully",
	}

	c.logEvent(w, l)
}

func (c *Config) logEvent(w http.ResponseWriter, l LogPayload) {
	//-----1. create json byteslice to send to auth service-------
	bSlice, err := json.MarshalIndent(l, "", "\t")
	c.handleErr(err)
	buffer := bytes.NewBuffer(bSlice) // reader format

	log.Println("Sending to logger service - buffer", buffer)

	// ----2. call the auth service from broker service ----
	respFromAuth, err := http.Post(
		"http://logger-service/log",
		"application/json",
		buffer,
	)
	if err != nil {
		log.Println("err while sending POST request to logger service from broker service")
		c.errorJSON(w, err)
		return
	}
	defer respFromAuth.Body.Close()

	log.Println("Response received from logger service:",
		respFromAuth.Status)

	//------3. make sure to receive correct status code--------
	if respFromAuth.StatusCode == http.StatusUnauthorized {
		log.Println("http.StatusUnauthorized received from logger service in broker service")
		c.errorJSON(w, errors.New("invalid credentials"))
		return
	}

	if respFromAuth.StatusCode != http.StatusAccepted {
		log.Println("http.StatusAccepted NOT received from logger service in broker service:", respFromAuth.Status)
		c.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	// if response code is ok, store response body
	var jsonFromAuth JsonResponse
	err = json.NewDecoder(respFromAuth.Body).Decode(&jsonFromAuth)
	log.Println("respFromAuth", respFromAuth)
	if err != nil {
		log.Println("err in logger service, while decoding json - respFromAuth", err)
		c.errorJSON(w, err)
		return
	}
	// check for err flag inside json
	if jsonFromAuth.ErrorPresent {
		log.Println("User is unauthorized", jsonFromAuth)
		c.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	//-----4. user is authenticated
	log.Println("User is authorized", jsonFromAuth)
	pLoad := JsonResponse{
		ErrorPresent: false,
		Message:      "event logged successfully!",
		Data:         jsonFromAuth.Data,
	}

	c.writeJSON(w, http.StatusAccepted, pLoad)
}
