package assert

import (
	"testing"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/stretchr/testify/assert"
)

// AssertFalse dummy struct
type AssertNil struct {
}

func init() {
	function.Register(&AssertNil{})
}

// Name of function
func (s *AssertNil) Name() string {
	return "assertNil"
}

// Sig - function signature
func (s *AssertNil) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeAny}, false
}

// Eval - function implementation
func (s *AssertNil) Eval(params ...interface{}) (interface{}, error) {
	actual := params[0]

	t := testing.T{}
	return assert.Nil(&t, actual), nil
}
