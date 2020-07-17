package update

import (
	azstorage "git.tibco.com/git/product/ipaas/wi-azstorage.git/src/app/Azurestorage/connector/connection"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support/connection"
)

//Input structure
type Input struct {
	Connection connection.Manager     `md:"Connection,required"`
	Operation  string                 `md:"operation,required"`
	Service    string                 `md:"service,required"`
	Insertflag bool                   `md:"insertFlag"`
	Input      map[string]interface{} `md:"input,required"`
}

//Output structure
type Output struct {
	Output interface{} `md:"output"` //

}

//FromMap method
func (i *Input) FromMap(values map[string]interface{}) error {
	var err error
	i.Service, _ = coerce.ToString(values["service"])
	if err != nil {
		return err
	}
	i.Operation, _ = coerce.ToString(values["operation"])
	if err != nil {
		return err
	}
	i.Connection, _ = azstorage.GetSharedConfiguration(values["Connection"])
	if err != nil {
		return err
	}
	i.Insertflag, _ = coerce.ToBool(values["insertFlag"])
	if err != nil {
		return err
	}
	i.Input, _ = coerce.ToObject(values["input"])
	if err != nil {
		return err
	}

	return nil
}

//ToMap method
func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"operation":  i.Operation,
		"service":    i.Service,
		"Connection": i.Connection,
		"insertFlag": i.Insertflag,
		"input":      i.Input,
	}
}

//ToMap Output
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"output": o.Output,
	}
}

//FromMap Output
func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	o.Output, _ = (values["output"])
	if err != nil {
		return err
	}

	return nil
}
