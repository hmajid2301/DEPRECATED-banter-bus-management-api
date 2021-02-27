package database

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
	*mongo.Database
	Client   *mongo.Client
	Logger   *log.Logger
	Host     string
	Port     int
	Username string
	Password string
	Name     string
	MaxConns int
	Timeout  int
}

// NewMongoDB sets up the Database
func NewMongoDB(
	logger *log.Logger,
	host string,
	port int,
	username string,
	password string,
	name string,
	maxConns int,
	timeout int,
) (*MongoDB, error) {
	db := &MongoDB{
		Logger:   logger,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Name:     name,
		MaxConns: maxConns,
		Timeout:  timeout,
	}

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
	db.Client = client
	db.Database = client.Database(name)
	return db, nil
}

// CloseDB closes the database
func (db *MongoDB) CloseDB() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout))
	defer cancel()

	err := db.Client.Disconnect(ctx)
	if err != nil {
		db.Logger.Errorf("Failed to disconnect from database, %s.", err)
	}
}

// Ping is used to check if the database is still connected to the app.
func (db *MongoDB) Ping() bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	err := db.Client.Ping(ctx, readpref.Primary())
	return err == nil
}

// Insert adds a new entry to the database.
func (db *MongoDB) Insert(collectionName string, document Document) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	db.Logger.WithFields(log.Fields{
		"collection": collectionName,
		"document":   document,
	}).Debug("Inserting document into database.")
	collection := db.Collection(collectionName)

	ok, err := collection.InsertOne(ctx, document)
	if err != nil {
		db.Logger.Error(err, ok)
		return false, err
	}

	var inserted = true
	if ok.InsertedID == nil {
		db.Logger.Error("No elements inserted.")
		inserted = false
	}

	return inserted, nil
}

// InsertMultiple adds multiple entries to the database at once.
func (db *MongoDB) InsertMultiple(collectionName string, documents Documents) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	db.Logger.WithFields(log.Fields{
		"collection": collectionName,
		"documents":  documents,
	}).Debug("Inserting multiple documents into database.")

	inserts := documents.ToInterface()
	collection := db.Collection(collectionName)
	_, err := collection.InsertMany(ctx, inserts)
	if err != nil {
		db.Logger.Error(err)
		return err
	}

	return nil
}

// Get retrieves an entry from the database.
func (db *MongoDB) Get(collectionName string, filter map[string]string, document Document) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	db.Logger.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
		"document":   document,
	}).Debug("Getting document from database.")
	collection := db.Collection(collectionName)
	err := collection.FindOne(ctx, filter).Decode(document)
	return err
}

// GetAll entries from the database.
func (db *MongoDB) GetAll(collectionName string, documents Documents) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	db.Logger.WithFields(log.Fields{
		"collection": collectionName,
		"documents":  documents,
	}).Debug("Getting all documents from database.")
	collection := db.Collection(collectionName)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		db.Logger.Errorf("failed to get objects: %v", err)
		return err
	}

	err = cursor.All(ctx, documents)
	if err != nil {
		db.Logger.Errorf("failed to transform object: %v", err)
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
	subDocuments SubDocuments,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	db.Logger.WithFields(log.Fields{
		"collection":   collectionName,
		"filter":       filter,
		"subDocuments": subDocuments,
	}).Debug("Getting sub-object from database.")
	collection := db.Collection(collectionName)

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

	err = aggregate.All(ctx, subDocuments)
	return err
}

// Delete removes a document from the database.
func (db *MongoDB) Delete(collectionName string, filter map[string]string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	db.Logger.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
	}).Debug("Deleting document from database.")
	collection := db.Collection(collectionName)

	ok, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		db.Logger.Error(err)
		return false, err
	}

	var deleted = true
	if ok.DeletedCount == 0 {
		db.Logger.Error("No elements deleted.")
		deleted = false
	}

	return deleted, nil
}

// RemoveCollection removes a collection from the database.
func (db *MongoDB) RemoveCollection(collectionName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	db.Logger.WithFields(log.Fields{
		"collection": collectionName,
	}).Warn("Deleting collection from database.")
	err := db.Collection(collectionName).Drop(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Update updates an existing "top-level" object in the database i.e. a user or a game.
func (db *MongoDB) Update(
	collectionName string,
	filter map[string]string,
	document Document,
) (bool, error) {
	db.Logger.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
		"document":   document,
	}).Debug("Updating document in the database.")

	updated, err := db.modifyEntry(collectionName, filter, document, "$set")
	return updated, err
}

// UpdateObject updates an existing "sub" object in the database i.e. a question inside a game.
// It is used to add fields to an object.
func (db *MongoDB) UpdateObject(
	collectionName string,
	filter map[string]string,
	update map[string]interface{},
) (bool, error) {
	db.Logger.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
		"update":     update,
	}).Debug("Updating sub-object in the database.")

	updated, err := db.modifyEntry(collectionName, filter, update, "$set")
	return updated, err
}

// RemoveObject updates an existing sub object in the database i.e. a question inside a game.
func (db *MongoDB) RemoveObject(
	collectionName string,
	filter map[string]string,
	remove map[string]interface{},
) (bool, error) {
	db.Logger.WithFields(log.Fields{
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
	subDocument NewSubDocument,
) (bool, error) {
	db.Logger.WithFields(log.Fields{
		"collection":  collectionName,
		"filter":      filter,
		"subDocument": subDocument,
	}).Debug("Adding subDocument to existing entry in the database.")

	updated, err := db.modifyEntry(collectionName, filter, subDocument, "$push")
	return updated, err
}

// RemoveFromList appends an new entry to an array in the database.
func (db *MongoDB) RemoveFromList(
	collectionName string,
	filter map[string]string,
	subDocument SubDocument,
) (bool, error) {
	db.Logger.WithFields(log.Fields{
		"collection":  collectionName,
		"filter":      filter,
		"subDocument": subDocument,
	}).Debug("Removing item from existing entry in the database.")

	updated, err := db.modifyEntry(collectionName, filter, subDocument, "$pull")
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

	collection := db.Collection(collectionName)

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
