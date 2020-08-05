package assert

import (
	"fmt"
	"reflect"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

// AssertTrue dummy struct
type AssertTrue struct {
}

func init() {
	function.Register(&AssertTrue{})
}

// Name of function
func (s *AssertTrue) Name() string {
	return "assertTrue"
}

// Sig - function signature
func (s *AssertTrue) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeBool}, false
}

// Eval - function implementation
func (s *AssertTrue) Eval(params ...interface{}) (interface{}, error) {
	actual := params[0]
	at, ok := actual.(bool)
	//log.RootLogger().Debugf("Start AssertEqual function with expected: %+v, actual: %+v\n", true, actual)

	if !ok {
		return false, fmt.Errorf("different data type, expected: boolean, actual: %s\n", reflect.TypeOf(actual).String())
	}

	if !at {
		return false, fmt.Errorf("different data value, expected: true, actual: %v\n", actual)
	}

	return true, nil
}
