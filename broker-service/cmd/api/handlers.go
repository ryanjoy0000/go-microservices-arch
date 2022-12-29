package main

import (
	"net/http"
)

func (c Config) BrokerHandler(w http.ResponseWriter, r *http.Request) {
	// set a json response
	payload := jsonResponse{
		ErrorPresent: false,
		Message:      "Hit the broker service",
	}

	err := c.writeJSON(w, http.StatusAccepted, payload)
	c.handleErr(err)
}
