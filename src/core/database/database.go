// Package database ...
package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	log "github.com/sirupsen/logrus"
)

// Config is data required to connect to the database.
type Config struct {
	Host         string
	Port         string
	DatabaseName string
	Username     string
	Password     string
}

var _client *mongo.Client
var _database *mongo.Database
var _ctx = context.TODO()

// InitialiseDatabase initializes the connection to the database.
func InitialiseDatabase(config Config) {
	log.Info("Connecting to database.")

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/", config.Username, config.Password, config.Host, config.Port)
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(_ctx, clientOptions)

	if err != nil {
		log.Error("Failed to connect to database.", err)
		panic(err)
	}

	log.Info("Connected to database.")

	_client = client
	_database = client.Database(config.DatabaseName)
}

// Ping is used to check if the database is still connected to the app.
func Ping() bool {
	err := _client.Ping(_ctx, readpref.Primary())
	if err != nil {
		log.Error("Failed to ping database.", err)
		return false
	}

	return true
}

// Insert adds a new entry to the database.
func Insert(collectionName string, object interface{}) (bool, error) {
	log.WithFields(log.Fields{
		"collection": collectionName,
		"object":     object,
	}).Debug("Inserting object into database.")
	collection := _database.Collection(collectionName)

	ok, err := collection.InsertOne(_ctx, object)
	if err != nil {
		log.Error(err, ok)
		return false, err
	}

	var inserted = true
	if ok.InsertedID == nil {
		log.Error("No elements inserted.")
		inserted = false
	}

	return inserted, nil
}

// InsertMultiple adds multiple entries to the database at once.
func InsertMultiple(collectionName string, object []interface{}) error {
	log.WithFields(log.Fields{
		"collection": collectionName,
		"object":     object,
	}).Debug("Inserting multiple objects into database.")
	collection := _database.Collection(collectionName)

	_, err := collection.InsertMany(_ctx, object)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// Get retrieves an entry from the database.
func Get(collectionName string, filter interface{}, model interface{}) error {
	log.WithFields(log.Fields{
		"collection": collectionName,
		"filer":      filter,
		"model":      model,
	}).Debug("Getting object from database.")
	collection := _database.Collection(collectionName)
	encodedFilter, marshalErr := bson.Marshal(filter)
	if marshalErr != nil {
		return marshalErr
	}

	err := collection.FindOne(_ctx, encodedFilter).Decode(model)
	return err
}

// GetAll entries from the database.
func GetAll(collectionName string, model interface{}) error {
	log.WithFields(log.Fields{
		"collection": collectionName,
		"model":      model,
	}).Debug("Getting multiple objects from database.")
	collection := _database.Collection(collectionName)

	cursor, err := collection.Find(_ctx, bson.M{})
	if err != nil {
		log.Error("Failed to get object", err)
	}

	if err = cursor.All(_ctx, model); err != nil {
		log.Error("Failed to transform object", err)
	}

	return err
}

// Delete removes an entry from the database.
func Delete(collectionName string, filter interface{}) (bool, error) {
	log.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
	}).Debug("Deleting object from database.")
	collection := _database.Collection(collectionName)

	ok, err := collection.DeleteOne(_ctx, filter)
	if err != nil {
		log.Error(err)
		return false, err
	}

	var deleted = true
	if ok.DeletedCount == 0 {
		log.Error("No elements deleted.")
		deleted = false
	}

	return deleted, nil
}

// RemoveCollection removes a collection from the database.
func RemoveCollection(collectionName string) {
	log.WithFields(log.Fields{
		"collection": collectionName,
	}).Warn("Deleting collection from database.")
	err := _database.Collection(collectionName).Drop(_ctx)
	if err != nil {
		panic(err)
	}
}

// UpdateEntry updates an existing entry in the database.
func UpdateEntry(collectionName string, filter interface{}, update interface{}) (bool, error) {
	log.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
		"update":     update,
	}).Debug("Update item in database.")
	collection := _database.Collection(collectionName)

	update = bson.M{"$set": update}
	ok, err := collection.UpdateOne(_ctx, filter, update)
	if err != nil {
		log.Error("Failed to update entry", err)
		return false, err
	}

	var updated = false
	if ok.ModifiedCount > 0 {
		updated = true
	}
	return updated, nil
}

// AppendToEntry appends an new entry to an array in the database.
func AppendToEntry(collectionName string, filter interface{}, add interface{}) (bool, error) {
	log.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
		"add":        add,
	}).Debug("Adding item to existing entry in database.")
	collection := _database.Collection(collectionName)

	ok, err := collection.UpdateOne(_ctx, filter, bson.M{"$push": add})
	if err != nil {
		log.Error("Failed to append to entry", err)
		return false, err
	}

	var updated = false
	if ok.ModifiedCount > 0 {
		updated = true
	}
	return updated, nil
}
