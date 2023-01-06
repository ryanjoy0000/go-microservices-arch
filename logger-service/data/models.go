package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbTimeout      = time.Second * 5
	dbName         = "logsDB"
	collectionName = "logsCollection"
)

var mongoClient *mongo.Client

// Single Log Entry
type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// Model is a type for the data package
type Models struct {
	LogEntry LogEntry
}

// Creates an instance of the data package. Returns Model struct which has all types available to the app
func New(mClient *mongo.Client) Models {
	// set the given db handle
	mongoClient = mClient

	// create an instance of model and return
	l1 := LogEntry{}
	m1 := Models{
		LogEntry: l1,
	}
	return m1
}

// -----------------------------------------------------
// DB methods

// insert a log entry
func (l *LogEntry) Insert(inputLog LogEntry) error {
	// set the DB and the collection
	collection := mongoClient.Database(dbName).Collection(collectionName)

	// modify input log entry for insertion
	l1 := LogEntry{
		Name:      inputLog.Name,
		Data:      inputLog.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// insert into collection
	_, err := collection.InsertOne(context.TODO(), l1)
	if err != nil {
		log.Println("error while inserting log into collection:", err)
		return err
	}
	return nil

}

// get all logs from collection
func (l *LogEntry) GetAllLogs() ([]*LogEntry, error) {

	// set the DB and the collection
	collection := mongoClient.Database(dbName).Collection(collectionName)

	//create a find options instance
	findOptions := options.Find()

	// set the value for sort field
	findOptions.SetSort(bson.D{{"created_at", -1}})

	// get matching records from collection
	mongoCursor, err := collection.Find(context.TODO(), bson.D{}, findOptions)
	if err != nil {
		log.Println("collection find error:", err)
		return nil, err
	}

	// set timeout context
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	defer mongoCursor.Close(ctx)

	// define list
	var logs []*LogEntry

	// scan through fetched logs
	for mongoCursor.Next(ctx) {
		var l1 LogEntry
		err := mongoCursor.Decode(&l1)
		if err != nil {
			log.Println("error decoding log:", err)
			return nil, err
		} else {
			// add to logs list
			logs = append(logs, &l1)
		}
	}

	return logs, nil
}

func (l *LogEntry) GetALog(id string) (*LogEntry, error) {

	// set the DB and the collection
	collection := mongoClient.Database(dbName).Collection(collectionName)

	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("error in getting object id:", err)
		return nil, err
	}

	var l1 LogEntry

	// set timeout context
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	// find record and store in l1
	err = collection.FindOne(ctx, bson.M{"_id": docId}).Decode(&l1)
	if err != nil {
		log.Println("error in finding single log:", err)
		return nil, err
	}

	return &l1, nil
}

func (l *LogEntry) DropCollection(id string) error {
	// set timeout context
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	// set the DB and the collection
	collection := mongoClient.Database(dbName).Collection(collectionName)

	err := collection.Drop(ctx)
	if err != nil {
		log.Println("error in dropping collection", err)
		return err
	}
	return nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	// set timeout context
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	// set the DB and the collection
	collection := mongoClient.Database(dbName).Collection(collectionName)

	docId, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		log.Println("error in getting doc id:", err)
		return nil, err
	}

	filterBsonObj := bson.M{"_id": docId}
	updObj := bson.D{
		{"$set", bson.D{
			{"name", l.Name},
			{"data", l.Data},
			{"updated_at", time.Now()},
		}},
	}

	updResult, err := collection.UpdateOne(
		ctx,
		filterBsonObj,
		updObj,
	)
	if err != nil {
		log.Println("error in updating doc:", err)
		return nil, err
	}

	return updResult, nil
}
