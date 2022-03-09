# pmongo
Pmongo is a wrapper on go mongo driver: go.mongodb.org/mongo-driver/mongo/options customized to easily create and manage mongo connection(s).

# Depedency
go get github.com/phil-inc/pmongo

# usage
Provide DBConfig or Slice of DBConfig with connection URL and Db name in system properties.
Example code to init mongo connection:
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

Example code for using the connection created:
```
GetConnection().Save(context, structure)
```
context: refer to go context
the structure is any go structure.

# Unit test
```bash
make test
```

# fmt
```bash
make fmt
```
