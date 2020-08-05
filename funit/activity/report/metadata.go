package report

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings of the activity
type Settings struct {
}

// Input of the activity
type Input struct {
}

// Output of the activity
type Output struct {
	Output interface{} `md:"output"`
}

// ToMap converts activity input to a map
func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{}
}

// FromMap sets activity input values from a map
func (i *Input) FromMap(values map[string]interface{}) error {

	return nil
}

// ToMap converts activity output to a map
func (o *Output) ToMap() map[string]interface{} {

	return map[string]interface{}{
		"output": o.Output,
	}
}

// FromMap sets activity output values from a map
func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	if o.Output, err = coerce.ToObject(values["output"]); err != nil {
		return err
	}

	return nil
}
