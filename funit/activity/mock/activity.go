package mock

import (
	"github.com/project-flogo/core/activity"
	logger "github.com/project-flogo/core/support/log"
)

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})
var log = logger.ChildLogger(logger.RootLogger(), "funit.activity.mock")

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Activity is a stub for creating a contract
type Activity struct {
	MockValue interface{}
}

// New creates a new Activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	settings := ctx.Settings()
	mockvalue, set := settings["mockOutput"]
	log.Infof("mockvalue: %v\n", mockvalue)
	if set {
		return &Activity{MockValue: mockvalue}, nil
	} else {
		return &Activity{}, nil
	}
}

// Metadata implements activity.Activity.Metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements activity.Activity.Eval
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	if a.MockValue != nil {
		ctx.SetOutput("output", a.MockValue)
	}
	
	return true, nil
}