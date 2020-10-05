// Package database ...
package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
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
var ctx = context.TODO()

func InitialiseDatabase(config DatabaseConfig) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/", config.Username, config.Password, config.Host, config.Port)
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("cat")
	database = client.Database(config.DatabaseName)
}

func Insert(collectionName string, id string, object interface{}) {
	collection := database.Collection(collectionName)
	itemToAdd := bson.M{id: object}
	_, err := collection.InsertOne(ctx, itemToAdd)
	if err != nil {
		log.Fatal(err)
	}
}

func Get(collectionName string, id string) *mongo.SingleResult {
	collection := database.Collection(collectionName)
	itemToGet := bson.M{"id": id}
	item := collection.FindOne(ctx, itemToGet)
	return item
}

func Delete(collectionName string, id string) {
	collection := database.Collection(collectionName)
	itemToDelete := bson.M{"id": id}
	_, err := collection.DeleteOne(ctx, itemToDelete)
	if err != nil {
		log.Fatal(err)
	}
}

func PartialUpdate(collectionName string, itemUpdate interface{}, id string) {
	collection := database.Collection(collectionName)
	filter := bson.M{"id": id}
	update := bson.M{"$set": itemUpdate}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}
}

func Disconnect() {
	err := client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
