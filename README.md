# pmongo
Pmongo is wrapper of mongo_driver of go.mongodb.org/mongo-driver/mongo/options customized as per phil's
requirement.

go get github.com/phil-inc/pmongo

Provide DBConfig or Slice of DBConfig with connection url and Db name in system properties.
Example code to init mongo connectino:
```
func InitMongoClient() {
	host := getOrDefault("db.host", "localhost:21707").(string)
	dbName := getOrDefault("db.name", "someDB").(string)
	dbConfig := core.DBConfig{HostURL: host, DBName: dbName}
	core.SetupMongoDB(dbConfig)
	mongoAppDB = core.ConnectionByName(dbName)
}
```

Example code to use mongo connection:
```
enrollCore.GetConnection().Save(context, structure)
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
