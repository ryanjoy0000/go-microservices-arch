package main

import (
	"log"
	"net/http"

	"github.com/ryanjoy0000/go-microservices-arch/logger-service/data"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (c *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// set a json response
	var reqPayload JSONPayload

	// store json to data
	err := c.readJSON(w, r, &reqPayload)
	if err != nil {
		log.Println("err while reading json into data", err)
		c.errorJSON(w, err, http.StatusBadRequest)

	}

	// create log entry
	l1 := data.LogEntry{
		Name: reqPayload.Name,
		Data: reqPayload.Data,
	}

	// insert log entry
	err = c.Models.LogEntry.Insert(l1)
	if err != nil {
		log.Println("err while inserting log", err)
		c.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	resp := JsonResponse{
		ErrorPresent: false,
		Message:      "logged into DB",
	}

	c.writeJSON(w, http.StatusAccepted, resp)
}
