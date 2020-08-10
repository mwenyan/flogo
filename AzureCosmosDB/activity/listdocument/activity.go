package listdocument

import (
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
var log = logger.ChildLogger(logger.RootLogger(), "cosmosdb.activity.listdocument")

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
	token := common.GetAuthToken(cnx.MasterKey, today, "get", "docs", fmt.Sprintf("dbs/%s/colls/%s", db, coll))

	req, err := http.NewRequest("get", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("x-ms-date", today)
	req.Header.Set("x-ms-version", cnx.APIVersion)

	if input.Data["x-ms-activity-id"] != nil {
		req.Header.Set("x-ms-activity-id", input.Data["x-ms-activity-id"].(string))
	}

	if input.Data["x-ms-session-token"] != nil {
		req.Header.Set("x-ms-session-token", input.Data["x-ms-session-token"].(string))
	}

	if input.Data["x-ms-continuation"] != nil {
		req.Header.Set("x-ms-continuation", input.Data["x-ms-continuation"].(string))
	}

	if input.Data["x-ms-max-item-count"] != nil {
		req.Header.Set("x-ms-max-item-count", strconv.Itoa(input.Data["x-ms-continuation"].(int)))
	}

	if input.Data["x-ms-consistency-level"] != nil {
		req.Header.Set("x-ms-consistency-level", input.Data["x-ms-consistency-level"].(string))
	}

	if input.Data["If-None-Match"] != nil {
		req.Header.Set("If-None-Match", input.Data["If-None-Match"].(string))
	}

	if input.Data["x-ms-documentdb-partitionkeyrangeid"] != nil {
		req.Header.Set("x-ms-documentdb-partitionkeyrangeid", input.Data["x-ms-documentdb-partitionkeyrangeid"].(string))
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
	if resp.Header.Get("x-ms-continuation") != "" {
		output["x-ms-continuation"] = resp.Header.Get("x-ms-continuation")
	}

	if result["code"] != nil {
		output["code"] = result["code"]
		output["message"] = result["message"]
	} else {
		for k, v := range result {
			output[k] = v
		}
	}
	ctx.SetOutput("output", output)
	return true, nil
}
