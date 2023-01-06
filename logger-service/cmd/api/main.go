package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ryanjoy0000/go-microservices-arch/logger-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort    = ":80"
	rpcPort    = ":5001"
	mongoUrl   = "mongodb://mongo:27017"
	grpcPort   = ":50001"
	maxTimeout = 15 * time.Second
)

var mongoClient *mongo.Client

type Config struct {
	Models data.Models
}

func main() {

	var app Config

	// connect to Mongo DB
	mongoClient, err := connectMongoDB()
	if err != nil {
		log.Panic(err)
	}

	// check for successul DB connection
	if mongoClient != nil {
		log.Println("Connected to Mongo DB")
		// set up config
		app = Config{
			Models: data.New(mongoClient),
		}

	}

	// context for DB
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, maxTimeout)
	// defer cancel func
	defer cancelFunc()

	// close DB connection after use
	defer closeMongoDB(ctx)

	// set up web server
	log.Println("Starting auth service on port", webPort)

	// define http server
	srv := &http.Server{
		Addr:    webPort,
		Handler: app.routes(),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic("error starting service", err)
	}

}

func connectMongoDB() (*mongo.Client, error) {
	// set the mongo DB connection options
	clientOptions := options.Client().ApplyURI(mongoUrl)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	mongoClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting mongo DB", err)
		return nil, err
	}

	return mongoClient, err
}

func closeMongoDB(ctx context.Context) {
	err := mongoClient.Disconnect(ctx)
	if err != nil {
		log.Panic(err)
	}
}
