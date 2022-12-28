package main

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Error   bool        `json:"errorPresent"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (c Config) BrokerHandler(w http.ResponseWriter, r *http.Request) {
	// set a json response
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker service",
	}

	// encode to json - byte slice
	jsonByteSlice, err := json.MarshalIndent(payload, "", "\t")
	c.handleErr(err)

	// set http headers to json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonByteSlice)
}
