package core

import (
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var dbConnection Db

func SecondaryConnection() *DBConnection {
	dbOptions := options.Database().SetReadPreference(readpref.SecondaryPreferred())
	return &DBConnection{DB: dbConnection.Client.Database(dbConnection.Config.DBName, dbOptions)}
}

func Connection() *DBConnection {
	return &DBConnection{DB: dbConnection.Client.Database(dbConnection.Config.DBName)}
}

func SetupMongoDB() error {
	config := DBConfig{HostURL: os.Getenv("db.main.url"), DBName: "MONGO_DB_NAME"}
	if err := Setup(config); err != nil {
		log.Fatalf("Error setting up data layer : %s %+v.\n", err, config)
		return err
	}

	newPdb, err := Get("DB_URL")
	if err != nil {
		log.Fatalf("DB connection error : %s.\n", err)
	}

	dbConnection = newPdb

	return nil
}
