package createdocument

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"net/http"

	"github.com/TIBCOSoftware/azure/AzureCosmosDB/common"
	"github.com/TIBCOSoftware/azure/AzureCosmosDB/connector/cosmosdb"

	"github.com/project-flogo/core/activity"
	logger "github.com/project-flogo/core/support/log"
)

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})
var log = logger.ChildLogger(logger.RootLogger(), "cosmosdb.activity.createdocument")

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Activity is a stub for creating a contract
type Activity struct {
}

// New creates a new Activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	return &Activity{}, nil
}

// Metadata implements activity.Activity.Metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements activity.Activity.Eval
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return false, err
	}

	cnx := input.Connection.GetConnection().(*cosmosdb.CosmosDBSharedConfigManager).GetClientConfiguration()

	today := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	db := input.Data["database"].(string)
	coll := input.Data["collection"].(string)
	url := fmt.Sprintf("%s/dbs/%s/colls/%s/docs", common.GetHost(cnx.Account), db, coll)
	token := common.GetAuthToken(cnx.MasterKey, today, "post", "docs", fmt.Sprintf("dbs/%s/colls/%s", db, coll))
	doc := input.Data["document"].(map[string]interface{})
	docid := input.Data["id"]
	if docid != nil {
		doc["id"] = docid
	}
	body, err := json.Marshal(doc)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("post", url, bytes.NewBuffer(body))
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("x-ms-date", today)
	req.Header.Set("x-ms-version", cnx.APIVersion)
	if input.Data["x-ms-documentdb-is-upsert"] != nil {
		req.Header.Set("x-ms-documentdb-is-upsert", strconv.FormatBool(input.Data["x-ms-documentdb-is-upsert"].(bool)))
	}

	if input.Data["x-ms-activity-id"] != nil {
		req.Header.Set("x-ms-activity-id", input.Data["x-ms-activity-id"].(string))
	}

	if input.Data["x-ms-documentdb-partitionkey"] != nil {
		pkeys := make([]string, 0)
		pkeys = append(pkeys, input.Data["x-ms-documentdb-partitionkey"].(string))
		pkeysbytes, err := json.Marshal(pkeys)
		if err != nil {
			return false, err
		}
		req.Header.Set("x-ms-documentdb-partitionkey", string(pkeysbytes))
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{Timeout: time.Duration(cnx.Timeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(respbody, &result)
	if err != nil {
		return false, err
	}

	output := make(map[string]interface{})
	output["x-ms-activity-id"] = resp.Header.Get("x-ms-activity-id")
	output["x-ms-session-token"] = resp.Header.Get("x-ms-session-token")
	if result["code"] != nil {
		output["code"] = result["code"]
		output["message"] = result["message"]
	} else {
		output["id"] = result["id"]
	}
	ctx.SetOutput("output", output)
	return true, nil
}
