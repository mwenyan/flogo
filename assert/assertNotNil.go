package assert

import (
	"testing"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/stretchr/testify/assert"
)

// AssertFalse dummy struct
type AssertNotNil struct {
}

func init() {
	function.Register(&AssertNotNil{})
}

// Name of function
func (s *AssertNotNil) Name() string {
	return "assertNotNil"
}

// Sig - function signature
func (s *AssertNotNil) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeAny}, false
}

// Eval - function implementation
func (s *AssertNotNil) Eval(params ...interface{}) (interface{}, error) {
	actual := params[0]

	t := testing.T{}
	return assert.NotNil(&t, actual), nil
}
