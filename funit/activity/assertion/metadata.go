package assert

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings of the activity
type Settings struct {
}

// Input of the activity
type Input struct {
	Name         string                 `md:"name"`
	ErrorMessage string                 `md:"msg"`
	AbortOnFail  bool                   `md:"abort"`
	Data         map[string]interface{} `md:"input"`
}

// Output of the activity
type Output struct {
	Output interface{} `md:"output"`
}

// ToMap converts activity input to a map
func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"input": i.Data,
		"abort": i.AbortOnFail,
		"name":  i.Name,
		"msg":   i.ErrorMessage,
	}
}

// FromMap sets activity input values from a map
func (i *Input) FromMap(values map[string]interface{}) error {
	var err error
	if i.Data, err = coerce.ToObject(values["input"]); err != nil {
		return err
	}
	if i.AbortOnFail, err = coerce.ToBool(values["abort"]); err != nil {
		return err
	}
	if i.Name, err = coerce.ToString(values["name"]); err != nil {
		return err
	}
	if i.ErrorMessage, err = coerce.ToString(values["msg"]); err != nil {
		return err
	}
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
