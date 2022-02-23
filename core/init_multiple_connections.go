package core

import (
	"log"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var dbConnectionMap map[string]Db

func SecondaryConnectionFromMap(dbName string) *DBConnection {
	dbOptions := options.Database().SetReadPreference(readpref.SecondaryPreferred())
	return &DBConnection{DB: dbConnectionMap[dbName].Client.Database(dbConnection.Config.DBName, dbOptions)}
}

func Connections(dbName string) *DBConnection {
	return &DBConnection{DB: dbConnectionMap[dbName].Client.Database(dbConnection.Config.DBName)}
}

func SetupMongoDBs(configs []DBConfig) error {
	for _, config := range configs {
		if err := Setup(config); err != nil {
			log.Fatalf("Error setting up data layer : %s %+v.\n", err, config)
			return err
		}

		newPdb, err := Get(config.DBName)
		if err != nil {
			log.Fatalf("DB connection error : %s.\n", err)
		}
		dbConnectionMap[config.DBName] = newPdb
	}

	return nil
}
