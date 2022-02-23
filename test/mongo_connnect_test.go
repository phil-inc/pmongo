package test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/phil-inc/pmongo/core"
	"github.com/stretchr/testify/assert"
)

func TestConnectivity(t *testing.T) {
	ctx := context.Background()
	os.Setenv("db.main.url", "mongodb://mongo:mongo@mongodb:27017")
	os.Setenv("db.main.name", "root-db")
	config := core.DBConfig{HostURL: "mongodb://mongo:mongo@mongodb:27017", DBName: "root-db"}
	core.SetupMongoDB(config)
	assertcommon(core.Connection() != nil, true, t, "TestConnectivity")

	data_setup(ctx, core.Connection(), "is", "mfr")

	q := core.Q{
		"internalStatus": "is",
	}
	osl := new(OrderStatusMappingLookUpInfo)
	core.Connection().Find(ctx, q, osl)
	assertcommon(osl.MfrStatus, "mfr", t, "TestConnectivity")
	dispose(ctx, core.Connection())
}

func TestConnectivitySecondaryConnectivityTest(t *testing.T) {
	ctx := context.Background()
	os.Setenv("db.main.url", "mongodb://mongo:mongo@mongodb:27017")
	os.Setenv("db.main.name", "root-db")
	config := core.DBConfig{HostURL: "mongodb://mongo:mongo@mongodb:27017", DBName: "root-db"}
	core.SetupMongoDB(config)
	assertcommon(core.Connection() != nil, true, t, "TestConnectivity")

	data_setup(ctx, core.Connection(), "is", "mfr")

	q := core.Q{
		"internalStatus": "is",
	}
	osl := new(OrderStatusMappingLookUpInfo)
	core.SecondaryConnection().Find(ctx, q, osl)
	assertcommon(osl.MfrStatus, "mfr", t, "TestConnectivity")
	dispose(ctx, core.Connection())

}

func TestMultiConnectivity(t *testing.T) {
	ctx := context.Background()
	config1 := core.DBConfig{HostURL: "mongodb://mongo:mongo@mongodb:27017", DBName: "root-db"}
	config2 := core.DBConfig{HostURL: "mongodb://mongo:mongo@mongodb:27017", DBName: "card-db"}
	core.SetupMongoDB(config1, config2)
	assertcommon(core.ConnectionByName("root-db") != nil, true, t, "TestMultiConnectivity")
	assertcommon(core.ConnectionByName("card-db") != nil, true, t, "TestMultiConnectivity")

	data_setup(ctx, core.ConnectionByName("root-db"), "is", "mfr")

	q := core.Q{
		"internalStatus": "is",
	}
	osl := new(OrderStatusMappingLookUpInfo)
	core.ConnectionByName("root-db").Find(ctx, q, osl)
	assertcommon(osl.MfrStatus, "mfr", t, "TestMultiConnectivity")

	data_setup(ctx, core.ConnectionByName("card-db"), "is2", "mfr2")

	q = core.Q{
		"internalStatus": "is2",
	}
	osl2 := new(OrderStatusMappingLookUpInfo)
	core.ConnectionByName("card-db").Find(ctx, q, osl2)
	assertcommon(osl2.MfrStatus, "mfr2", t, "TestMultiConnectivity")

	osl = new(OrderStatusMappingLookUpInfo)
	core.ConnectionByName("root-db").Find(ctx, core.Q{"internalStatus": "is2"}, osl)
	assertcommon(osl.MfrStatus, "", t, "TestMultiConnectivity")

	osl = new(OrderStatusMappingLookUpInfo)
	core.ConnectionByName("card-db").Find(ctx, core.Q{"internalStatus": "is"}, osl)
	assertcommon(osl.MfrStatus, "", t, "TestMultiConnectivity")
	dispose(ctx, core.ConnectionByName("root-db"))
	dispose(ctx, core.ConnectionByName("card-db"))
}

func TestMultiSecondaryDbConnectivity(t *testing.T) {
	ctx := context.Background()
	config1 := core.DBConfig{HostURL: "mongodb://mongo:mongo@mongodb:27017", DBName: "root-db"}
	config2 := core.DBConfig{HostURL: "mongodb://mongo:mongo@mongodb:27017", DBName: "card-db"}
	core.SetupMongoDB(config1, config2)
	data_setup(ctx, core.ConnectionByName("root-db"), "is", "mfr")

	q := core.Q{
		"internalStatus": "is",
	}
	osl := new(OrderStatusMappingLookUpInfo)
	core.SecondaryConnectionByName("root-db").Find(ctx, q, osl)
	assertcommon(osl.MfrStatus, "mfr", t, "TestMultiSecondaryDbConnectivity")

	data_setup(ctx, core.ConnectionByName("card-db"), "is2", "mfr2")

	q = core.Q{
		"internalStatus": "is2",
	}
	osl2 := new(OrderStatusMappingLookUpInfo)
	core.SecondaryConnectionByName("card-db").Find(ctx, q, osl2)
	assertcommon(osl2.MfrStatus, "mfr2", t, "TestMultiSecondaryDbConnectivity")

	osl = new(OrderStatusMappingLookUpInfo)
	core.SecondaryConnectionByName("root-db").Find(ctx, core.Q{"internalStatus": "is2"}, osl)
	assertcommon(osl.MfrStatus, "", t, "TestMultiSecondaryDbConnectivity")

	osl = new(OrderStatusMappingLookUpInfo)
	core.SecondaryConnectionByName("card-db").Find(ctx, core.Q{"internalStatus": "is"}, osl)
	assertcommon(osl.MfrStatus, "", t, "TestMultiSecondaryDbConnectivity")
	dispose(ctx, core.SecondaryConnectionByName("root-db"))
	dispose(ctx, core.SecondaryConnectionByName("card-db"))
}

func assertcommon(actual interface{}, expected interface{}, t *testing.T, testCase string) {
	result := assert.Equal(t, expected, actual, fmt.Sprintf("Error: Test failed mismatch %s != %s", actual, expected))
	fmt.Println(result)
	if !result {
		t.Fatal(errors.New(fmt.Sprintf("%s failed", testCase)))
	} else {
		fmt.Printf(testCase)
	}
}

func data_setup(ctx context.Context, conn *core.DBConnection, is, mfr string) {
	o := OrderStatusMappingLookUpInfo{InternalStatus: is,
		MfrStatus: mfr,
	}
	conn.Save(ctx, o)
}

func dispose(ctx context.Context, conn *core.DBConnection) {
	osl := new(OrderStatusMappingLookUpInfo)
	conn.RemoveAll(ctx, core.Q{}, osl)
}

type OrderStatusMappingLookUpInfo struct {
	InternalStatus string `json:"internalStatus" bson:"internalStatus" binding:"required"`
	MfrStatus      string `json:"mfrStatus" bson:"mfrStatus" binding:"required"`
}

// CollectionName
func (_ OrderStatusMappingLookUpInfo) CollectionName() string {
	return "orderStatusMappingLookUpInfo"
}
