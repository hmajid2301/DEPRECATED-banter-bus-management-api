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

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
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
func (db *MongoDB) Insert(collectionName string, objectToInsert models.Document) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	db.logger.WithFields(log.Fields{
		"collection": collectionName,
		"object":     objectToInsert,
	}).Debug("Inserting object into database.")
	collection := db.database.Collection(collectionName)

	ok, err := collection.InsertOne(ctx, objectToInsert)
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
func (db *MongoDB) InsertMultiple(collectionName string, object models.Documents) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	db.logger.WithFields(log.Fields{
		"collection": collectionName,
		"object":     object,
	}).Debug("Inserting multiple objects into database.")

	objectsToInsert := object.ToInterface()
	collection := db.database.Collection(collectionName)
	_, err := collection.InsertMany(ctx, objectsToInsert)
	if err != nil {
		db.logger.Error(err)
		return err
	}

	return nil
}

// Get retrieves an entry from the database.
func (db *MongoDB) Get(collectionName string, filter map[string]string, model models.Document) error {
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
func (db *MongoDB) GetAll(collectionName string, model models.Documents) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	db.logger.WithFields(log.Fields{
		"collection": collectionName,
		"model":      model,
	}).Debug("Getting all objects from database.")
	collection := db.database.Collection(collectionName)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		db.logger.Errorf("failed to get objects: %v", err)
		return err
	}

	err = cursor.All(ctx, model)
	if err != nil {
		db.logger.Errorf("failed to transform object: %v", err)
	}

	return err
}

// GetSubObject retrieves a subdocument from the database.
// This method should be used to retrieve a single subdocument from an array element.
// The filter is what to use to filter to that subdocument i.e. `{"username": "virat_kohli"}`.
// The parentField being the first field you want to see i.e. `question_pools`.
// The condition is which of the array elements to get from the parentField i.e. `{$$this.pool_name", "my_pool"}`.
// The model must be a slice/array where to unmarshal the BSON to i.e. `[]models.QuestionPool`.
// This method will return the first element the matches the condition.
func (db *MongoDB) GetSubObject(
	collectionName string,
	filter map[string]string,
	parentField string,
	condition []string,
	model models.SubDocuments,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	db.logger.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
		"model":      model,
	}).Debug("Getting sub-object from database.")
	collection := db.database.Collection(collectionName)

	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: filter},
		},
		{
			{
				Key: "$replaceRoot",
				Value: bson.M{
					"newRoot": bson.M{
						"$arrayElemAt": []interface{}{
							bson.M{
								"$filter": bson.M{
									"input": fmt.Sprintf("$%s", parentField),
									"cond": bson.M{
										"$eq": condition,
									},
								},
							},
							0,
						},
					},
				},
			},
		},
	}

	aggregate, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}

	err = aggregate.All(ctx, model)
	return err
}

// Delete removes an entry from the database.
func (db *MongoDB) Delete(collectionName string, filter map[string]string) (bool, error) {
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

// Update updates an existing "top-level" object in the database i.e. a user or a game.
func (db *MongoDB) Update(
	collectionName string,
	filter map[string]string,
	objectToUpdate models.Document,
) (bool, error) {
	db.logger.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
		"update":     objectToUpdate,
	}).Debug("Updating item in the database.")

	updated, err := db.modifyEntry(collectionName, filter, objectToUpdate, "$set")
	return updated, err
}

// UpdateObject updates an existing "sub" object in the database i.e. a question inside a game.
// It is used to add fields to an object.
func (db *MongoDB) UpdateObject(
	collectionName string,
	filter map[string]string,
	objectToAdd map[string]interface{},
) (bool, error) {
	db.logger.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
		"update":     objectToAdd,
	}).Debug("Adding new item to object in the database.")

	updated, err := db.modifyEntry(collectionName, filter, objectToAdd, "$set")
	return updated, err
}

// RemoveObject updates an existing sub object in the database i.e. a question inside a game.
func (db *MongoDB) RemoveObject(
	collectionName string,
	filter map[string]string,
	remove map[string]interface{},
) (bool, error) {
	db.logger.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
		"update":     remove,
	}).Debug("Removing object in the database.")

	updated, err := db.modifyEntry(collectionName, filter, remove, "$unset")
	return updated, err
}

// AppendToList appends an new entry to an array in the database.
func (db *MongoDB) AppendToList(
	collectionName string,
	filter map[string]string,
	add models.NewSubDocument,
) (bool, error) {
	db.logger.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
		"add":        add,
	}).Debug("Adding item to existing entry in the database.")

	updated, err := db.modifyEntry(collectionName, filter, add, "$push")
	return updated, err
}

// RemoveFromList appends an new entry to an array in the database.
func (db *MongoDB) RemoveFromList(
	collectionName string,
	filter map[string]string,
	remove models.SubDocument,
) (bool, error) {
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
	filter map[string]string,
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
