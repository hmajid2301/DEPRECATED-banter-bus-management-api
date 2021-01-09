package repository

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDB is a representation of the MongoDB server to connect to
type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
	logger   *log.Logger

	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
	MaxConns     int
	Timeout      int
}

// NewMongoDB sets up the Database
func NewMongoDB(
	logger *log.Logger,
	host string,
	port int,
	username string,
	password string,
	databaseName string,
	maxConns int,
	timeout int,
) (*MongoDB, error) {
	db := &MongoDB{
		logger:       logger,
		Host:         host,
		Port:         port,
		Username:     username,
		Password:     password,
		DatabaseName: databaseName,
		MaxConns:     maxConns,
		Timeout:      timeout}

	logger.Info("Connecting to database.")

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/", username, password, host, port)

	logger.Debugf("Database connection string: %s", uri)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	clientOptions = clientOptions.SetMaxPoolSize(uint64(maxConns))
	clientOptions = clientOptions.SetConnectTimeout(time.Duration(timeout) * time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return &MongoDB{}, fmt.Errorf("error while connecting to database: %w", err)
	}

	logger.Info("Connected to database.")
	db.client = client
	db.database = client.Database(databaseName)
	return db, nil
}

// CloseDB closes the database
func (db *MongoDB) CloseDB() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout))
	defer cancel()

	err := db.client.Disconnect(ctx)
	if err != nil {
		db.logger.Errorf("Failed to disconnect from database, %s.", err)
	}
}

// Ping is used to check if the database is still connected to the app.
func (db *MongoDB) Ping() bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	err := db.client.Ping(ctx, readpref.Primary())
	return err == nil
}

// Insert adds a new entry to the database.
func (db *MongoDB) Insert(collectionName string, object interface{}) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	db.logger.WithFields(log.Fields{
		"collection": collectionName,
		"object":     object,
	}).Debug("Inserting object into database.")
	collection := db.database.Collection(collectionName)

	ok, err := collection.InsertOne(ctx, object)
	if err != nil {
		db.logger.Error(err, ok)
		return false, err
	}

	var inserted = true
	if ok.InsertedID == nil {
		db.logger.Error("No elements inserted.")
		inserted = false
	}

	return inserted, nil
}

// InsertMultiple adds multiple entries to the database at once.
func (db *MongoDB) InsertMultiple(collectionName string, object []interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	db.logger.WithFields(log.Fields{
		"collection": collectionName,
		"object":     object,
	}).Debug("Inserting multiple objects into database.")

	collection := db.database.Collection(collectionName)
	_, err := collection.InsertMany(ctx, object)
	if err != nil {
		db.logger.Error(err)
		fmt.Println(err)
		return err
	}

	return nil
}

// Get retrieves an entry from the database.
func (db *MongoDB) Get(collectionName string, filter interface{}, model interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	db.logger.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
		"model":      model,
	}).Debug("Getting object from database.")
	collection := db.database.Collection(collectionName)
	encodedFilter, marshalErr := bson.Marshal(filter)
	if marshalErr != nil {
		return marshalErr
	}

	err := collection.FindOne(ctx, encodedFilter).Decode(model)
	return err
}

// GetAll entries from the database.
func (db *MongoDB) GetAll(collectionName string, model interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	db.logger.WithFields(log.Fields{
		"collection": collectionName,
		"model":      model,
	}).Debug("Getting multiple objects from database.")
	collection := db.database.Collection(collectionName)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		db.logger.Errorf("Failed to get object: %v.", err)
	}

	if err = cursor.All(ctx, model); err != nil {
		db.logger.Errorf("Failed to transform object: %v.", err)
	}

	return err
}

// Delete removes an entry from the database.
func (db *MongoDB) Delete(collectionName string, filter interface{}) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	db.logger.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
	}).Debug("Deleting object from database.")
	collection := db.database.Collection(collectionName)

	ok, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		db.logger.Error(err)
		return false, err
	}

	var deleted = true
	if ok.DeletedCount == 0 {
		db.logger.Error("No elements deleted.")
		deleted = false
	}

	return deleted, nil
}

// RemoveCollection removes a collection from the database.
func (db *MongoDB) RemoveCollection(collectionName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	db.logger.WithFields(log.Fields{
		"collection": collectionName,
	}).Warn("Deleting collection from database.")
	err := db.database.Collection(collectionName).Drop(ctx)
	if err != nil {
		return err
	}

	return nil
}

// UpdateEntry updates an existing entry in the database.
func (db *MongoDB) UpdateEntry(collectionName string, filter interface{}, update interface{}) (bool, error) {
	db.logger.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
		"update":     update,
	}).Debug("Update item in the database.")

	updated, err := db.modifyEntry(collectionName, filter, update, "$set")
	return updated, err
}

// RemoveEntry updates an existing entry in the database.
func (db *MongoDB) RemoveEntry(collectionName string, filter interface{}, update interface{}) (bool, error) {
	db.logger.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
		"update":     update,
	}).Debug("Updating item in the database.")

	updated, err := db.modifyEntry(collectionName, filter, update, "$unset")
	return updated, err
}

// AppendToEntry appends an new entry to an array in the database.
func (db *MongoDB) AppendToEntry(collectionName string, filter interface{}, add interface{}) (bool, error) {
	db.logger.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
		"add":        add,
	}).Debug("Adding item to existing entry in the database.")

	updated, err := db.modifyEntry(collectionName, filter, add, "$push")
	return updated, err
}

// RemoveFromEntry appends an new entry to an array in the database.
func (db *MongoDB) RemoveFromEntry(collectionName string, filter interface{}, remove interface{}) (bool, error) {
	db.logger.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
		"remove":     remove,
	}).Debug("Removing item from existing entry in the database.")

	updated, err := db.modifyEntry(collectionName, filter, remove, "$pull")
	return updated, err
}

func (db *MongoDB) modifyEntry(
	collectionName string,
	filter interface{},
	modify interface{},
	operation string,
) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	collection := db.database.Collection(collectionName)

	ok, err := collection.UpdateOne(ctx, filter, bson.M{operation: modify})
	if err != nil {
		return false, err
	}

	var updated = false
	if ok.ModifiedCount > 0 {
		updated = true
	}

	return updated, nil
}
