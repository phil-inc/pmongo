package core

import (
	"log"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var dbConnectionMap map[string]Db = make(map[string]Db)

// Since we all pool of connection per db in db map.
// use this this method if needed secondary connection for specific db name.
func SecondaryConnectionByName(dbName string) *DBConnection {
	dbOptions := options.Database().SetReadPreference(readpref.SecondaryPreferred())
	return &DBConnection{DB: dbConnectionMap[dbName].Client.Database(dbConnectionMap[dbName].Config.DBName, dbOptions)}
}

//Will return the only db secondary connection.
// if your app is connected to multiple mongodb, then please use SecondaryConnectionByName
// this will return nil
func SecondaryConnection() *DBConnection {
	if len(dbConnectionMap) == 1 {
		keys := reflect.ValueOf(dbConnectionMap).MapKeys()
		dbConnection := dbConnectionMap[keys[0].String()]
		dbOptions := options.Database().SetReadPreference(readpref.SecondaryPreferred())
		return &DBConnection{DB: dbConnection.Client.Database(dbConnection.Config.DBName, dbOptions)}
	}
	return nil
}

//Will return connection by db name from connection map.
func ConnectionByName(dbName string) *DBConnection {
	return &DBConnection{DB: dbConnectionMap[dbName].Client.Database(dbConnectionMap[dbName].Config.DBName)}
}

// Will return the only connection by name or if your app is connected to single mongodb.
// if your app is connected to multiple mongodb, then please use ConnectionByName this will return nil
func Connection() *DBConnection {
	if len(dbConnectionMap) == 1 {
		keys := reflect.ValueOf(dbConnectionMap).MapKeys()
		dbConnection := dbConnectionMap[keys[0].String()]
		return &DBConnection{DB: dbConnection.Client.Database(dbConnection.Config.DBName)}
	}
	return nil
}

//Set up mongodb with provided vargs of configs.
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
