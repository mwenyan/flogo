package assert

import (
	"fmt"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/app"
	logger "github.com/project-flogo/core/support/log"
)

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})
var log = logger.ChildLogger(logger.RootLogger(), "funit.activity.assertion")

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

	output := make(map[string]interface{})

	isTrue := input.Data["assertion"].(bool)
	output["assertResult"] = isTrue

	msgOutput := "OK"
	if !isTrue {
		msg := input.ErrorMessage
		msgvars := make([]interface{}, 0)
		if msg != "" {
			for k, v := range input.Data {
				if k != "assertion" {
					msgvars = append(msgvars, v)
				}
			}
			msgOutput = fmt.Sprintf(msg, msgvars...)
		} else {
			msgOutput = "Failed"
		}

		if input.AbortOnFail {
			return false, fmt.Errorf(msgOutput)
		}
	}

	ctx.SetOutput("output", output)

	report, ok := app.GetValue("funitReport")
	if !ok {
		report = make(map[string]string)
	}

	report.(map[string]string)[input.Name] = msgOutput
	app.SetValue("funitReport", report)
	return true, nil
}
