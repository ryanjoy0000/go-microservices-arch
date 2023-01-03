package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type JsonResponse struct {
	ErrorPresent bool        `json:"errorPresent"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data,omitempty"`
}

// Common Functions

// convert received single JSON (inside request body) to data
func (c *Config) readJSON_OLD(w http.ResponseWriter, r *http.Request, data interface{}) error {
	// maxBytes := int64(1024 * 1024) // 1MB

	// Limiting the size of incoming request body
	// r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	// reads received JSON (from req body) & converts to data
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	// ------------------------>
	// DEBUG
	var bs []byte
	r.Body.Read(bs)
	body := string(bs)

	b, _ := ioutil.ReadAll(r.Body)
	bs2 := string(b)

	log.Println("Decode process - r.Body", r.Body)
	log.Println("Decode process - r.Body", body)
	log.Println("Decode process - r.Body", bs2)

	log.Println("Decode process - err", err)
	log.Println("Decode process - data", data)
	// ------------------------>
	if err != nil {
		log.Println("Err while decoding json to data:", err)
		return err
	}

	// check if there are more than 1 JSON file
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("request body has more than one JSON value")
	}

	return nil
}

func (c *Config) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Err while reading request body:", err)
		return err
	}

	log.Println("reading request body", string(body))

	defer r.Body.Close()

	err = json.Unmarshal(body, data)
	if err != nil {
		log.Println("Err while unmarshalling json to data:", err)
		return err
	}

	log.Println("request body is unmarshalled to data:", data)

	return nil
}

// convert data to JSON and send as response
func (c *Config) writeJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	var err error = nil

	// convert data to JSON
	bSlice, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// check if headers are present
	if len(headers) > 0 {
		for key, val := range headers[0] {
			// set the header key and val
			w.Header()[key] = val
		}
	}

	// set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// set the http status
	w.WriteHeader(status)

	// send the data on response
	_, err = w.Write(bSlice)
	if err != nil {
		return err
	}

	return err
}

// write err msg as JSON to response
func (c *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {

	// set err status code
	errCode := http.StatusBadRequest
	if len(status) > 0 {
		errCode = status[0]
	}

	// set payload
	var errPayLoad JsonResponse
	errPayLoad.ErrorPresent = true
	errPayLoad.Message = err.Error()

	// send
	return c.writeJSON(w, errCode, errPayLoad)
}
