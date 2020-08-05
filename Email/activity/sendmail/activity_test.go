package sendmail

import (
	"testing"

	"io/ioutil"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/stretchr/testify/assert"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestRegistered(t *testing.T) {
	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}

func TestRecipients(t *testing.T) {
	re := "123@tibco .com ,234@tibc.com"
	assert.Equal(t, "123@tibco .com,234@tibc.com", getRecipentsStr(re))
	assert.Equal(t, []string{"123@tibco .com", "234@tibc.com"}, getRecipents(re))

	re2 := "123@tibco .com "
	assert.Equal(t, "123@tibco .com", getRecipentsStr(re2))
	assert.Equal(t, []string{"123@tibco .com"}, getRecipents(re2))
}
