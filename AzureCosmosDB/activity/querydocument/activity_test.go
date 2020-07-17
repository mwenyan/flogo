package querydocument

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/TIBCOSoftware/azure/AzureCosmosDB/connector/cosmosdb"

	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/resolve"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestQueryDocWithoutParams(t *testing.T) {

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
	query := "SELECT * FROM abc f"
	input["database"] = "tempdb"
	input["collection"] = "tempcoll"
	input["x-ms-documentdb-query-enablecrosspartition"] = false

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

	err = tc.SetInputObject(&Input{Data: input, Connection: mgr, Query: query})
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

func TestQueryDocWithParams(t *testing.T) {

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
	//	inputmp := make(map[string]interface{})
	//	inputmp["Value"] = `{"observers":[],"issuer":"Alice","amount":"10000","currency":"USD","owner":"Alice"}`
	//	v, err := json.Marshal(inputmp)
	//	complexObject := &data.ComplexObject{}
	//	err = json.Unmarshal([]byte(v), complexObject)

	query := "SELECT * FROM abc f WHERE f.id = @familyId"
	v := `{"@familyId":"WakefieldFamily"}`
	input := make(map[string]interface{})
	var doc interface{}
	input["database"] = "tempdb"
	input["collection"] = "tempcoll"
	input["x-ms-documentdb-query-enablecrosspartition"] = true
	err = json.Unmarshal([]byte(v), &doc)
	if err != nil {
		fmt.Errorf("json parse error: %v", err)
	}

	input["parameters"] = doc
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

	err = tc.SetInputObject(&Input{Data: input, Connection: mgr, Query: query})
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
