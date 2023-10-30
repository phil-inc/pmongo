package core

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseData represents the basic fields that can be embedded in MongoDB documents.
type BaseData struct {
	ID          interface{} `json:"id" bson:"_id,omitempty"`
	CreatedDate *time.Time  `json:"createdDate" bson:"createdDate,omitempty"`
	UpdatedDate *time.Time  `json:"updatedDate" bson:"updatedDate,omitempty"`
}

// DBRef represents a MongoDB database reference.
type DBRef struct {
	Collection string      `bson:"$ref"`
	Id         interface{} `bson:"$id"`
}

// NewDBRef creates a new MongoDB database reference (DBRef) for a given collection name and object ID.
func NewDBRef(collectionName string, ID interface{}) *DBRef {
	return &DBRef{Collection: collectionName, Id: ObjectID(ID)}
}

// NewObjectID generates a new MongoDB ObjectID.
func NewObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

/*
Functions handling Mongo Driver ObjectID
*/
//ObjectID returns objectID from interface
func ObjectID(id interface{}) primitive.ObjectID {
	if id != nil {
		switch v := id.(type) {
		case string:
			i, _ := primitive.ObjectIDFromHex(v)
			return i
		case primitive.ObjectID:
			return v
		}
	}
	return [12]byte{}
}

// RefID extracts the hexadecimal representation of an ObjectID from a DBRef, if present.
func RefID(dbRef *DBRef) string {
	if dbRef == nil {
		return ""
	}
	if bsonID, ok := dbRef.Id.(primitive.ObjectID); ok {
		return bsonID.Hex()
	}
	return ""
}

// StringID returns the hexadecimal string representation of an ID
func StringID(ID interface{}) string {
	if ID != nil {
		switch v := ID.(type) {
		case string:
			return v
		case primitive.ObjectID:
			return v.Hex()
		default:
			return ""
		}
	}
	return ""
}

// StringID returns the hexadecimal string representation of the ID field within a BaseData structure
func (data *BaseData) StringID() string {
	if data == nil {
		return ""
	}
	return StringID(data.ID)
}

// StringIDsToObjectIDs converts a slice of hexadecimal string IDs into a slice of primitive.ObjectID.
func StringIDsToObjectIDs(stringIDs []string) []primitive.ObjectID {
	objectIDs := make([]primitive.ObjectID, 0)
	for _, stringID := range stringIDs {
		if stringID == "" {
			continue
		}
		objectIDs = append(objectIDs, ObjectID(stringID))
	}

	return objectIDs
}
