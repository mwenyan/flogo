package assert

import (
	"testing"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
	logger "github.com/project-flogo/core/support/log"
	"github.com/stretchr/testify/assert"
)

var log = logger.ChildLogger(logger.RootLogger(), "assert.assertEqual")

// AssertEqual dummy struct
type AssertEqual struct {
}

func init() {
	function.Register(&AssertEqual{})
}

// Name of function
func (s *AssertEqual) Name() string {
	return "assertEqual"
}

// Sig - function signature
func (s *AssertEqual) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeAny, data.TypeAny}, false
}

// Eval - function implementation
func (s *AssertEqual) Eval(params ...interface{}) (interface{}, error) {
	expected := params[0]
	actual := params[1]

	t := testing.T{}
	return assert.Equal(&t, expected, actual), nil
	/*
		et := reflect.TypeOf(expected)
		at := reflect.TypeOf(actual)

		if et != at {
			return false, nil
			//		return false, fmt.Errorf("different data type, expected: %s, actual: %s\n", et.String(), at.String())
		}

		if expected != actual {
			return false, nil
			//		return false, fmt.Errorf("not equal, expected: %v, actual: %v\n", expected, actual)
		}
	*/
	//	return true, nil
}
