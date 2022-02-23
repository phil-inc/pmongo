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
	core.SetupMongoDB()
	Assertcommon(core.Connection() != nil, true, t, "TestConnectivity")

	data_setup(ctx, core.Connection())

	q := core.Q{
		"internalStatus": "is",
	}
	osl := new(OrderStatusMappingLookUpInfo)
	core.Connection().Find(ctx, q, osl)
	Assertcommon(osl.MfrStatus, "mfr", t, "TestConnectivity")

}

func Assertcommon(actual interface{}, expected interface{}, t *testing.T, testCase string) {
	result := assert.Equal(t, expected, actual, fmt.Sprintf("Error: Test failed mismatch %s != %s", actual, expected))
	fmt.Println(result)
	if !result {
		t.Fatal(errors.New(fmt.Sprintf("%s failed", testCase)))
	}
}
func data_setup(ctx context.Context, conn *core.DBConnection) {
	o := OrderStatusMappingLookUpInfo{InternalStatus: "is",
		MfrStatus: "mfr",
	}
	conn.Save(ctx, o)
}

type OrderStatusMappingLookUpInfo struct {
	InternalStatus string `json:"internalStatus" bson:"internalStatus" binding:"required"`
	MfrStatus      string `json:"mfrStatus" bson:"mfrStatus" binding:"required"`
}

// CollectionName
func (_ OrderStatusMappingLookUpInfo) CollectionName() string {
	return "orderStatusMappingLookUpInfo"
}
