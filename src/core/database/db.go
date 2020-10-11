// Package database ...
package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/", config.Username, config.Password, config.Host, config.Port)
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	database = client.Database(config.DatabaseName)

	// defer func() {
	// 	if err = client.Disconnect(ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()
}

func Insert(collectionName string, object interface{}) error {
	collection := database.Collection(collectionName)
	_, err := collection.InsertOne(ctx, object)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func Get(collectionName string, filter interface{}, model interface{}) error {
	collection := database.Collection(collectionName)
	err := collection.FindOne(ctx, filter).Decode(&model)
	return err
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
