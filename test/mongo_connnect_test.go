package test

import (
	"os"
	"testing"

	"github.com/phil-inc/pmongo/core"
)

func TestSomeInsertion(t *testing.T) {
	os.Setenv("db.main.url", "mongodb://mongo:mongo@mongodb:27017")
	os.Setenv("db.main.name", "root-db")
	core.Connection()
}

type SomeObj struct {
	Value string `json:"value" bson:"value" binding:"required"`
}
