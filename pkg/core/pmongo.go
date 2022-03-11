package core

import (
	"context"
	"errors"
	"log"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Q query representation to hide bson.M type to single file
type Q map[string]interface{}

// connectionMap holds all the db connection per database name
var connectionMap = make(map[string]Db)

// Get creates new database connection
func Get(dbName string) (Db, error) {
	if db, ok := connectionMap[dbName]; ok {
		return db, nil
	}
	return Db{}, errors.New("database connection not available. Perform 'Setup' first")
}

// DB represents database connection which holds reference to global client and configuration for that database.
type Db struct {
	Config DBConfig
	Client *mongo.Client
}

type CursorOptions struct {
	BatchSize int32
	Limit     int64
	Skip      int64
	Sort      map[string]int
}

var ErrNoDocumentsFound = errors.New("mongo: no documents in result").Error()

// Setup the MongoDB connection based on passed in config. It can be called multiple times to setup connection to
// multiple MongoDB instances.
func Setup(dbConfig DBConfig) error {
	if dbConfig.HostURL == "" || dbConfig.DBName == "" {
		return errors.New("invalid connection info. Missing host and db info")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(dbConfig.HostURL))
	if err != nil {
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Printf("MongoDB %s connection failed : %s. Exiting the program.\n", dbConfig.DBName, err)
		return err
	}

	//starting with primary preferred, but individual query can change mode per copied session
	log.Printf("Connected to %s via pmongo successfully", dbConfig.DBName)

	/* Initialized database object with global database connection*/
	connectionMap[dbConfig.DBName] = Db{Config: dbConfig, Client: client}
	return nil
}

// DBConfig represents the configuration params needed for MongoDB connection
type DBConfig struct {
	HostURL, DBName string
	Mode            int
}

// Document interface implemented by structs that needs to be persisted. It should provide collection name,
// as in the database. Also, a way to create new object id before saving.
type Document interface {
	CollectionName() string
}

// DBConnection pmongo connection wrapper
type DBConnection struct {
	DB *mongo.Database
}

// collection returns a mgo.collection representation for given collection name and session
func (s *DBConnection) Collection(collectionName string) *mongo.Collection {
	return s.DB.Collection(collectionName)
}

// Save inserts the given document that represents the collection to the database.
func (s *DBConnection) Save(ctx context.Context, document Document) error {
	coll := s.Collection(document.CollectionName())
	_, err := coll.InsertOne(ctx, document)
	return err
}

// Update updates the given document based on given selector
func (s *DBConnection) Update(ctx context.Context, selector Q, document Document) error {
	coll := s.Collection(document.CollectionName())
	//_, err := coll.UpdateOne(ctx, selector, bson.M{"$set": document})
	_, err := coll.ReplaceOne(ctx, selector, document)
	return err
}

// Update updates the given document based on given selector
func (s *DBConnection) Upsert(ctx context.Context, selector Q, document Document) error {
	coll := s.Collection(document.CollectionName())
	//_, err := coll.UpdateOne(ctx, selector, bson.M{"$set": document})
	_, err := coll.ReplaceOne(ctx, selector, document, options.Replace().SetUpsert(true))
	return err
}

// UpdateByID updates the given document based on given id
func (s *DBConnection) UpdateByID(ctx context.Context, id string, result Document) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}
	return s.Update(ctx, Q{"_id": objID}, result)
}

// FindByID find the object by id. Returns error if it's not able to find the document. If document is found
// it's copied to the passed in result object.
func (s *DBConnection) FindByID(ctx context.Context, id string, result Document) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}
	return s.Find(ctx, Q{"_id": objID}, result)
}

// Find the data based on given query
func (s *DBConnection) Find(ctx context.Context, query Q, document Document) error {
	err := s.Collection(document.CollectionName()).FindOne(ctx, query).Decode(document)
	if err != nil {
		if err.Error() != mongo.ErrNoDocuments.Error() {
			log.Printf("Error fetching %s with query %s. Error: %s\n", document.CollectionName(), query, err)
		}
	}
	return err
}

// FindAll returns all the documents based on given query
func (s *DBConnection) FindAll(ctx context.Context, query Q, document Document) (interface{}, error) {
	curr, err := s.Collection(document.CollectionName()).Find(ctx, query)
	if err != nil {
		return nil, err
	}
	documents := slice(document)
	err = curr.All(ctx, documents)
	if err != nil {
		return nil, err
	}
	return results(documents)
}

func results(documents interface{}) (interface{}, error) {
	return reflect.ValueOf(documents).Elem().Interface(), nil
}

// slice returns the interface representation of actual collection type for returning list data
func slice(d Document) interface{} {
	documentType := reflect.TypeOf(d)
	documentSlice := reflect.MakeSlice(reflect.SliceOf(documentType), 0, 0)

	// Create a pointer to a slice value and set it to the slice
	return reflect.New(documentSlice.Type()).Interface()
}

// Exists check if the document exists for given query
func (s *DBConnection) Exists(ctx context.Context, query Q, document Document) (bool, error) {
	err := s.Find(ctx, query, document)
	if err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

//Remove removes the given document type based on the query
func (s *DBConnection) Remove(ctx context.Context, query Q, document Document) error {
	_, err := s.Collection(document.CollectionName()).DeleteOne(ctx, query)
	return err
}

// RemoveByID remove the object by id. Returns error if it's not able to find the document. If document is found
// it's copied to the passed in result object.
func (s *DBConnection) RemoveByID(ctx context.Context, id string, result Document) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}
	return s.Remove(ctx, Q{"_id": objID}, result)
}

//RemoveAll removes all the document matching given selector query
func (s *DBConnection) RemoveAll(ctx context.Context, query Q, document Document) error {
	_, err := s.Collection(document.CollectionName()).DeleteMany(ctx, query)
	return err
}

//RemoveAllWithCount removes all the document matching given selector query
func (s *DBConnection) RemoveAllWithCount(ctx context.Context, query Q, document Document) (int64, error) {
	res, err := s.Collection(document.CollectionName()).DeleteMany(ctx, query)
	if res == nil {
		return -1, err
	}
	return res.DeletedCount, err
}

//GetCursor gets a cursor to iterate over the documents returned by the selector
func (s *DBConnection) GetCursor(ctx context.Context, query Q, collectionName string, cursorOptions CursorOptions) (*mongo.Cursor, error) {
	opts := &options.FindOptions{
		BatchSize: &cursorOptions.BatchSize,
		Skip:      &cursorOptions.Skip,
		Sort:      cursorOptions.Sort,
		Limit:     &cursorOptions.Limit,
	}

	return s.Collection(collectionName).Find(ctx, query, opts)
}

//UpdateFieldValue updates the single field with a given value for a collection name based query
func (s *DBConnection) UpdateFieldValue(ctx context.Context, query Q, collectionName, field string, value interface{}) error {
	_, err := s.Collection(collectionName).UpdateOne(ctx, query, bson.M{"$set": bson.M{field: value}})

	return err
}
