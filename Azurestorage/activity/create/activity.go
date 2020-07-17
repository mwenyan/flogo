package create

import (
	"fmt"
	"io/ioutil"
	"regexp"

	commonutil "git.tibco.com/git/product/ipaas/wi-azstorage.git/src/app/Azurestorage"
	azstorage "git.tibco.com/git/product/ipaas/wi-azstorage.git/src/app/Azurestorage/connector/connection"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/log"
)

var re = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

const (
	ivConnection = "Connection"
	ivInput      = "input"
	ivService    = "service"
	ivOperation  = "operation"
	ovOutput     = "output"
	ovError      = "error"
)

var activityMd = activity.ToMetadata(&Input{}, &Output{})
var activityLog = log.ChildLogger(log.RootLogger(), "azure-storage-create")

func init() {
	err := activity.Register(&Activity{}, New)
	if err != nil {
		log.RootLogger().Error(err)
	}
}

// New functioncommon
func New(ctx1 activity.InitContext) (activity.Activity, error) {
	return &Activity{}, nil
}

// Activity is a stub for your Activity implementation
type Activity struct {
	operation    string
	service      string
	typeofUpload string
	path         string
	sastoken     string
	accountname  string
}

// Metadata implements activity.Activity.Metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

//Cleanup method
func (a *Activity) Cleanup() error {

	return nil
}

// Eval implements activity.Activity.Eval
func (a *Activity) Eval(context activity.Context) (done bool, err error) {
	activityLog.Debugf("Executing Activity [%s] ", context.Name())
	inputVal := &Input{}
	err = context.GetInputObject(inputVal)
	service := inputVal.Service
	if service == "" {
		return false, activity.NewError("service is not configured", "AZURE-STORAGE-1003", nil)
	}

	operation := inputVal.Operation
	if operation == "" {
		return false, activity.NewError("operation is not configured", "AZURE-STORAGE-1004", nil)
	}
	typeofUpload := "Byte Upload" //For this relese only byte upload is supported
	mcon, _ := inputVal.Connection.(*azstorage.AzStorageSharedConfigManager)
	if mcon.GetClientConfiguration().SAS == "" {
		activityLog.Error("Please re-authenticate your connection")
		return false, activity.NewError("SAS could not be generated", "AZURE-STORAGE-1005", nil)

	}
	sasToken := mcon.GetClientConfiguration().SAS
	accountName := mcon.GetClientConfiguration().Accountname

	paramMap := make(map[string]string)
	if inputVal != nil {
		inputMap := inputVal.Input
		if inputMap["parameters"] != nil {
			parameters := inputMap["parameters"]
			for k, v := range parameters.(map[string]interface{}) {
				paramMap[k] = fmt.Sprint(v)

			}
		}
	}
	azstorageAPIPath := commonutil.GetAzureStorageAPIpath(accountName, service, operation, paramMap)
	activityLog.Debug("Sending Request to Azure storage backend...")
	objectResponse := make(map[string]interface{})
	msgResponse := make(map[string]interface{})
	if service == "File" {
		resp, err := commonutil.DoCall("PUT", operation, azstorageAPIPath, paramMap, sasToken)

		if err != nil {
			activityLog.Error(err)
			return false, err
		}
		jsonResponseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			activityLog.Error(err)
			return false, fmt.Errorf("Error reading JSON response data after Create invocation, %s", err.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			err := commonutil.ErrorHandeler(resp, jsonResponseData)
			if err != nil {
				activityLog.Error(err)
				return false, err
			}
		}
		msgResponse["statusMessage"] = "Operation " + operation + " successfully completed."
		msgResponse["statusCode"] = resp.StatusCode
		msgResponse["isSuccess"] = true
		msgResponse["shareName"] = paramMap["shareName"]
		if operation == "Create Directory" {
			msgResponse["directoryPath"] = paramMap["directoryPath"]
			msgResponse["directoryName"] = paramMap["directoryName"]
		} else {
			msgResponse["directoryPath"] = paramMap["directoryPath"]
			msgResponse["fileName"] = paramMap["fileName"]
		}
		objectResponse[service] = msgResponse

	} else if service == "Blob" {
		if operation == "Upload Blob" {
			blobURL := azstorageAPIPath + "?" + sasToken[1:len(sasToken)]
			responseBlob, reserror := commonutil.UploadBlob(blobURL, typeofUpload, a.path, paramMap)
			if reserror != nil {
				return false, reserror
			}
			msgResponse["statusMessage"] = "Operation " + operation + " successfully completed."
			msgResponse["statusCode"] = responseBlob.StatusCode
			msgResponse["containerName"] = paramMap["containerName"]
			msgResponse["blobName"] = paramMap["blobName"]
			objectResponse[service] = msgResponse
		}
	}
	context.SetOutput("output", objectResponse)
	activityLog.Debugf("Execution of Activity [%s] " + context.Name() + " completed")
	return true, nil
}
