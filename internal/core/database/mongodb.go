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

	err = db.createIndex("question", "id", true)
	if err != nil {
		return &MongoDB{}, fmt.Errorf("error while creating index for question %w", err)
	}

	err = db.createIndex("story", "id", true)
	if err != nil {
		return &MongoDB{}, fmt.Errorf("error while creating index for story %w", err)
	}
	return db, nil
}

func (db *MongoDB) CloseDB() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout))
	defer cancel()

	err := db.Client.Disconnect(ctx)
	if err != nil {
		db.Logger.Errorf("Failed to disconnect from database, %s.", err)
	}
}

func (db *MongoDB) Ping() bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	err := db.Client.Ping(ctx, readpref.Primary())
	return err == nil
}

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

func (db *MongoDB) Get(collectionName string, filter map[string]interface{}, document Document) error {
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

func (db *MongoDB) GetRandom(
	collectionName string,
	filter map[string]interface{},
	limit int64,
	documents Documents,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	db.Logger.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
		"limit":      limit,
	}).Debug("Getting random fields from database.")

	collection := db.Collection(collectionName)
	pipeline := mongo.Pipeline{
		{
			{
				Key: "$match", Value: filter,
			},
		},
		{
			{
				Key: "$sample",
				Value: bson.M{
					"size": limit,
				},
			},
		},
	}

	aggregate, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}

	err = aggregate.All(ctx, documents)
	return err
}

func (db *MongoDB) GetUniqueValues(
	collectionName string,
	filter map[string]interface{},
	field string,
) ([]string, error) {
	db.Logger.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
		"field_name": field,
	}).Debug("Getting unique fields from database.")

	selector := fmt.Sprintf("$%s", field)
	pipeline := mongo.Pipeline{
		{
			{
				Key: "$match", Value: filter,
			},
		},
		{
			{
				Key: "$project",
				Value: bson.M{
					"item": 1,
					field: bson.M{
						"$ifNull": []string{selector, ""},
					},
				},
			},
		},
		{
			{
				Key: "$group",
				Value: bson.M{
					"_id":  "_id",
					"temp": bson.M{"$addToSet": selector},
				},
			},
		},
		{
			{
				Key: "$unset", Value: []string{"_id"},
			},
		},
	}

	unique, err := db.aggregate(collectionName, pipeline)
	return unique, err
}

func (db *MongoDB) GetUniqueKeys(
	collectionName string,
	filter map[string]interface{},
	fieldName string,
) ([]string, error) {
	db.Logger.WithFields(log.Fields{
		"collection": collectionName,
		"field_name": fieldName,
	}).Debug("Getting unique key from database.")

	selector := fmt.Sprintf("$%s", fieldName)
	pipeline := mongo.Pipeline{
		{
			{
				Key: "$match", Value: filter,
			},
		},
		{
			{
				Key: "$project", Value: bson.M{
					"_id": 0,
				},
			},
		},
		{
			{
				Key: "$project",
				Value: bson.M{
					"item": 1,
					"o": bson.M{
						"$objectToArray": selector,
					},
				},
			},
		},
		{
			{
				Key:   "$unwind",
				Value: "$o",
			},
		},
		{
			{
				Key: "$group",
				Value: bson.M{
					"_id":  "_id",
					"temp": bson.M{"$addToSet": "$o.k"},
				},
			},
		},
		{
			{
				Key: "$unset", Value: []string{"_id"},
			},
		},
	}

	unique, err := db.aggregate(collectionName, pipeline)
	return unique, err
}

func (db *MongoDB) aggregate(collectionName string, pipeline mongo.Pipeline) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	collection := db.Collection(collectionName)
	aggregate, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	model := []map[string][]string{}
	err = aggregate.All(ctx, &model)
	if err != nil {
		return []string{}, err
	} else if len(model) < 1 {
		return []string{}, nil
	}

	uniqueItems := model[0]["temp"]
	return uniqueItems, nil
}

func (db *MongoDB) GetAll(collectionName string, filter map[string]interface{}, documents Documents) error {
	options := options.Find()
	err := db.find(collectionName, filter, documents, options)
	return err
}

func (db *MongoDB) GetWithLimit(
	collectionName string,
	filter map[string]interface{},
	limit int64,
	documents Documents,
) error {
	options := options.Find()
	options.SetLimit(limit)
	err := db.find(collectionName, filter, documents, options)
	return err
}

func (db *MongoDB) find(
	collectionName string,
	filter map[string]interface{},
	documents Documents,
	options *options.FindOptions,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	db.Logger.WithFields(log.Fields{
		"collection": collectionName,
		"documents":  documents,
		"options":    options,
	}).Debug("Getting all documents from database.")
	collection := db.Collection(collectionName)

	cursor, err := collection.Find(ctx, filter, options)
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

func (db *MongoDB) Delete(collectionName string, filter map[string]interface{}) (bool, error) {
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

func (db *MongoDB) DeleteAll(collectionName string, filter map[string]interface{}) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	db.Logger.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
	}).Debug("Deleting documents from database.")
	collection := db.Collection(collectionName)

	ok, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		db.Logger.Error(err)
		return false, err
	}

	var deleted = true
	if ok.DeletedCount == 0 {
		db.Logger.Warning("No elements deleted, nothing matched filter.")
		deleted = false
	}

	return deleted, nil
}

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

func (db *MongoDB) Update(
	collectionName string,
	filter map[string]interface{},
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

func (db *MongoDB) UpdateObject(
	collectionName string,
	filter map[string]interface{},
	subDocument UpdateSubDocument,
) (bool, error) {
	db.Logger.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
		"document":   subDocument,
	}).Debug("Updating sub-object in the database.")

	updated, err := db.modifyEntry(collectionName, filter, subDocument, "$set")
	return updated, err
}

func (db *MongoDB) RemoveObject(
	collectionName string,
	filter map[string]interface{},
	subDocument UpdateSubDocument,
) (bool, error) {
	db.Logger.WithFields(log.Fields{
		"collection": collectionName,
		"filter":     filter,
		"document":   subDocument,
	}).Debug("Removing object in the database.")

	updated, err := db.modifyEntry(collectionName, filter, subDocument, "$unset")
	return updated, err
}

func (db *MongoDB) AppendToList(
	collectionName string,
	filter map[string]interface{},
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

func (db *MongoDB) RemoveFromList(
	collectionName string,
	filter map[string]interface{},
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
	filter map[string]interface{},
	modify interface{},
	operation string,
) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	collection := db.Collection(collectionName)

	update := bson.M{operation: modify}
	ok, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}

	var updated = false
	if ok.ModifiedCount > 0 {
		updated = true
	}

	return updated, nil
}

func (db *MongoDB) createIndex(collectionName string, field string, unique bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)
	defer cancel()

	mod := mongo.IndexModel{
		Keys:    bson.M{field: 1},
		Options: options.Index().SetUnique(unique),
	}

	collection := db.Collection(collectionName)
	_, err := collection.Indexes().CreateOne(ctx, mod)
	return err
}
