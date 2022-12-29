package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ryanjoy0000/go-microservices-arch/auth-service/data"
)

const webPort = ":80"
const maxDBConnAttempts = 10

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {

	var app Config

	// connect to DB
	dbConn := connectToDB()

	// check for successul DB connection
	if dbConn != nil {
		// set up config
		app = Config{
			DB:     dbConn,
			Models: data.New(dbConn),
		}
	}

	// set up web server
	log.Println("Starting auth service on port", webPort)

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

func connectToDB() *sql.DB {
	// get DB connection details from environment for security reasons
	dsn := os.Getenv("DSN")

	// try connecting DB until reaching max attempts
	for i := 0; i < maxDBConnAttempts; i++ {
		db, err := openDB(dsn)
		if err != nil {
			log.Println("DB not yet ready...Retry in 2 seconds")
		} else {
			log.Println("Connection to DB established...")
			return db
		}

		// delay each db connection attempt by 2 seconds
		time.Sleep(time.Second * 2)
	}

	log.Println("Aborting DB connection process...")
	return nil
}

func openDB(dsn string) (*sql.DB, error) {
	// open db connection
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// test db by ping
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
