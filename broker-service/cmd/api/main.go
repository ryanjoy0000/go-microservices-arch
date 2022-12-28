package main

import (
	"log"
	"net/http"
)

const webPort = ":80"

type Config struct{}

func main() {
	app := Config{}
	log.Println("Starting broker service on port", webPort)

	// define http server
	srv := &http.Server{
		Addr:    webPort,
		Handler: app.routes(),
	}

	// start the server
	err := srv.ListenAndServe()
	app.handleErr(err)
}

func (c Config) handleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
