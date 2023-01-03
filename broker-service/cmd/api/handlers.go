package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func (c *Config) BrokerHandler(w http.ResponseWriter, r *http.Request) {
	// set a json response
	payload := JsonResponse{
		ErrorPresent: false,
		Message:      "Hit the broker service",
	}

	err := c.writeJSON(w, http.StatusAccepted, payload)
	c.handleErr(err)
}

// listen request on FE, contact auth service, respond back to FE
func (c *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {

	log.Println("Received req in broker service", r.Method, r.URL)

	var reqPayload RequestPayload
	// store request from FE
	err := c.readJSON(w, r, &reqPayload)
	if err != nil {
		log.Println("Err while storing request from FE", err)
		c.errorJSON(w, err)
		return
	}

	// decide next steps based on action specified in request
	switch reqPayload.Action {
	case "auth":
		// authentication requested
		c.auth(w, reqPayload.Auth)
	default:
		c.errorJSON(w, errors.New("unknown request payload action"))
	}
}

func (c *Config) auth(w http.ResponseWriter, a AuthPayload) {
	log.Println("Starting authentication process as requested by FE") // <---------------------------------------

	//-----1. create json byteslice to send to auth service-------
	bSlice, err := json.MarshalIndent(a, "", "\t")
	c.handleErr(err)
	buffer := bytes.NewBuffer(bSlice) // reader format

	log.Println("Sending to auth service - buffer", buffer)

	// ----2. call the auth service from broker service ----
	respFromAuth, err := http.Post(
		"http://auth-service/auth",
		"application/json",
		buffer,
	)
	if err != nil {
		log.Println("err while sending POST request to auth service from broker service")
		c.errorJSON(w, err)
		return
	}
	defer respFromAuth.Body.Close()

	log.Println("Response received from auth service:",
		respFromAuth.Status)

	//------3. make sure to receive correct status code--------
	if respFromAuth.StatusCode == http.StatusUnauthorized {
		log.Println("http.StatusUnauthorized received from auth service in broker service")
		c.errorJSON(w, errors.New("invalid credentials"))
		return
	}

	if respFromAuth.StatusCode != http.StatusAccepted {
		log.Println("http.StatusAccepted NOT received from auth service in broker service:", respFromAuth.Status)
		c.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	// if response code is ok, store response body
	var jsonFromAuth JsonResponse
	err = json.NewDecoder(respFromAuth.Body).Decode(&jsonFromAuth)
	log.Println("respFromAuth", respFromAuth)
	if err != nil {
		log.Println("err in broker service, while decoding json - respFromAuth", err)
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
		Message:      "User is authenticated successfully!",
		Data:         jsonFromAuth.Data,
	}

	c.writeJSON(w, http.StatusAccepted, pLoad)
}
