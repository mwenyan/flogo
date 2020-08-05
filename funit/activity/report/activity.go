package report

import (
	"github.com/project-flogo/core/activity"
	logger "github.com/project-flogo/core/support/log"

	"github.com/project-flogo/core/app"
)

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})
var log = logger.ChildLogger(logger.RootLogger(), "funit.activity.report")

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Activity is a stub for creating a contract
type Activity struct {
}

type AssertResult struct {
	Name    string
	Pass    bool
	Message string
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

	output := make([]AssertResult, 0)
	report, ok := app.GetValue("funitReport")
	if ok {
		tests := report.(map[string]string)
		for k, v := range tests {
			test := AssertResult{}
			test.Name = k
			test.Message = v
			if v == "OK" {
				test.Pass = true
			} else {
				test.Pass = false
			}
			output = append(output, test)
		}
	}
	ctx.SetOutput("output", output)

	return true, nil
}
