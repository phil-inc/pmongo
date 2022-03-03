# pmongo
Pmongo is wrapper of mongo_driver of go.mongodb.org/mongo-driver/mongo/options customized as per phil's
requirement.

# Depedency
go get github.com/phil-inc/pmongo

# usage
Provide DBConfig or Slice of DBConfig with connection url and Db name in system properties.
Example code to init mongo connectino:
```
var mongoAppDB *core.DBConnection

func InitMongoClient() {
	host := getOrDefault("db.host", "localhost:21707").(string)
	dbName := getOrDefault("db.name", "someDB").(string)
	dbConfig := core.DBConfig{HostURL: host, DBName: dbName}
	core.SetupMongoDB(dbConfig)
	mongoAppDB = core.ConnectionByName(dbName)
}

func GetConnection() *core.DBConnection {
	return mongoAppDB
}
```

Example code to use mongo connection:
```
GetConnection().Save(context, structure)
```


# Unit test
```bash
make test
```

# fmt
```bash
make fmt
go fmt ./...
```
