package core

import (
	"log"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var dbConnectionMap map[string]Db = make(map[string]Db)

func SecondaryConnectionFromMap(dbName string) *DBConnection {
	dbOptions := options.Database().SetReadPreference(readpref.SecondaryPreferred())
	return &DBConnection{DB: dbConnectionMap[dbName].Client.Database(dbConnectionMap[dbName].Config.DBName, dbOptions)}
}

func SecondaryConnection() *DBConnection {
	if len(dbConnectionMap) == 1 {
		keys := reflect.ValueOf(dbConnectionMap).MapKeys()
		dbConnection := dbConnectionMap[keys[0].String()]
		dbOptions := options.Database().SetReadPreference(readpref.SecondaryPreferred())
		return &DBConnection{DB: dbConnection.Client.Database(dbConnection.Config.DBName, dbOptions)}
	}
	return nil
}

func ConnectionByName(dbName string) *DBConnection {
	return &DBConnection{DB: dbConnectionMap[dbName].Client.Database(dbConnectionMap[dbName].Config.DBName)}
}

func Connection() *DBConnection {
	if len(dbConnectionMap) == 1 {
		keys := reflect.ValueOf(dbConnectionMap).MapKeys()
		dbConnection := dbConnectionMap[keys[0].String()]
		return &DBConnection{DB: dbConnection.Client.Database(dbConnection.Config.DBName)}
	}
	return nil
}

func SetupMongoDB(configs ...DBConfig) error {
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
