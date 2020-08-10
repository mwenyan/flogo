package deletedocument

import (
	"fmt"
	"testing"

	"github.com/TIBCOSoftware/azure/AzureCosmosDB/connector/cosmosdb"

	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/resolve"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestDeleteDoc(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	mf := mapper.NewFactory(resolve.GetBasicResolver())
	iCtx := test.NewActivityInitContext(Settings{}, mf)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())

	input := make(map[string]interface{})
	input["database"] = "tempdb"
	input["collection"] = "tempcoll2"
	input["id"] = "mytest"
	input["x-ms-documentdb-partitionkey"] = "mytest"

	conn := make(map[string]interface{})
	conn["account"] = "jetbluecosmos"
	conn["key"] = "EPiTQHgKy5REhtu7v6vQN5HGeLMVHtMkL7078ZaUGVWsQlcIe5kDBmJCE73g22qSGyOKgQlGukPaAm6arBFebQ=="
	conn["api"] = "2018-12-31"
	conn["timeout"] = 10

	factory := &cosmosdb.CosmosDBConnectorFactory{}

	mgr, err := factory.NewManager(conn)
	if err != nil {
		fmt.Errorf("can't create connection manager: %v", err)
	}

	err = tc.SetInputObject(&Input{Data: input, Connection: mgr})
	if err != nil {
		fmt.Errorf("set input error: %v", err)
	}

	_, err = act.Eval(tc)
	if err != nil {
		fmt.Errorf("has error: %v", err)
	}

	//check output
	output := tc.GetOutput("output")
	fmt.Printf("output=%v\n", output)

}
