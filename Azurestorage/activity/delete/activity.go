package delete

import (
	"fmt"
	"io/ioutil"

	commonutil "git.tibco.com/git/product/ipaas/wi-azstorage.git/src/app/Azurestorage"
	azstorage "git.tibco.com/git/product/ipaas/wi-azstorage.git/src/app/Azurestorage/connector/connection"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/log"
)

const (
	ivConnection = "Connection"
	ivInput      = "input"
	ivService    = "service"
	ivOperation  = "operation"
	ovOutput     = "output"
	ovError      = "error"
)

// Activity is a stub for your Activity implementation
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
	operation   string
	service     string
	sastoken    string
	accountname string
}

var activityMd = activity.ToMetadata(&Input{}, &Output{})
var activityLog = log.ChildLogger(log.RootLogger(), "azure-storage-delete")

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
		resp, err := commonutil.DoCall("DELETE", operation, azstorageAPIPath, paramMap, sasToken)

		if err != nil {
			activityLog.Error(err)
			return false, err
		}
		activityLog.Debug("Received Response from azure storage backend")
		jsonResponseData, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			activityLog.Error(err)
			return false, fmt.Errorf("Error reading JSON response data after Create invocation, %s", err.Error())
		}
		defer resp.Body.Close()
		activityLog.Debugf("Response status for current operation is [%s]: ", resp.Status)

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
		if operation == "Delete Directory" {
			msgResponse["directoryPath"] = paramMap["directoryPath"]
			msgResponse["directoryName"] = paramMap["directoryName"]
		} else {
			msgResponse["directoryPath"] = paramMap["directoryPath"]
			msgResponse["fileName"] = paramMap["fileName"]
		}
		objectResponse[service] = msgResponse

	} else if service == "Blob" {
		if operation == "Delete Blob" {
			blobURL := azstorageAPIPath + "?" + sasToken[1:len(sasToken)]
			responseBlob, reserror := commonutil.DeleteBlob(blobURL, paramMap)
			if reserror != nil {
				return false, reserror
			}
			jsonResponseData, err := ioutil.ReadAll(responseBlob.Body)
			if err != nil {
				activityLog.Error(err)
				return false, err
			}
			msgResponse["statusMessage"] = "Operation " + operation + " successfully completed."
			msgResponse["statusCode"] = responseBlob.StatusCode
			msgResponse["responseBody"] = string(jsonResponseData)
			objectResponse[service] = msgResponse
		}
	}
	context.SetOutput("output", objectResponse)
	activityLog.Debugf("Execution of Activity [%s] " + context.Name() + " completed")
	return true, nil
}
