// Package database ...
package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	log "github.com/sirupsen/logrus"
)

//DatabaseConfig is ...
type DatabaseConfig struct {
	Host         string
	Port         string
	DatabaseName string
	Username     string
	Password     string
}

var client *mongo.Client
var database *mongo.Database
var collection *mongo.Collection
var ctx = context.TODO()

func InitialiseDatabase(config DatabaseConfig) {
	log.Info("Connecting to database.")
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/", config.Username, config.Password, config.Host, config.Port)
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Error("Failed to connect to database.", err)
		panic(err)
	}

	log.Info("Connected to database")
	database = client.Database(config.DatabaseName)
}

func Ping() bool {
	err := client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	return true
}

func Insert(collectionName string, object interface{}) error {
	log.WithFields(log.Fields{
		"collection": collectionName,
		"object":     object,
	}).Debug("Inserting object into database.")
	collection := database.Collection(collectionName)
	_, err := collection.InsertOne(ctx, object)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func InsertMultiple(collectionName string, object []interface{}) error {
	log.WithFields(log.Fields{
		"collection": collectionName,
		"object":     object,
	}).Debug("Inserting multiple objects into database.")
	collection := database.Collection(collectionName)
	_, err := collection.InsertMany(ctx, object)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func Get(collectionName string, filter interface{}, model interface{}) error {
	log.WithFields(log.Fields{
		"collection": collectionName,
		"filer":      filter,
		"model":      model,
	}).Debug("Getting object from database.")
	collection := database.Collection(collectionName)
	err := collection.FindOne(ctx, filter).Decode(model)
	return err
}

func GetAll(collectionName string, model interface{}) error {
	log.WithFields(log.Fields{
		"collection": collectionName,
		"model":      model,
	}).Debug("Getting multiple objects from database.")
	collection := database.Collection(collectionName)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Error("Failed to get object", err)
	}

	if err = cursor.All(ctx, model); err != nil {
		log.Error("Failed to transform object", err)
	}
	return err
}

func Delete(collectionName string, filter interface{}) {
	log.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
	}).Debug("Deleting object from database.")
	collection := database.Collection(collectionName)
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		panic(err)
	}
}

func RemoveCollection(collectionName string) {
	log.WithFields(log.Fields{
		"collection": collectionName,
	}).Warn("Deleting collection from database.")
	err := database.Collection(collectionName).Drop(ctx)
	if err != nil {
		panic(err)
	}
}

func UpdateItem(collectionName string, filter interface{}, update interface{}) {
	log.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
		"update":     update,
	}).Debug("Update item in database.")
	collection := database.Collection(collectionName)
	update = bson.M{"$set": update}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		panic(err)
	}
}
