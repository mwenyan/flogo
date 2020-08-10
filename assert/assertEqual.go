package assert

import (
	"testing"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/stretchr/testify/assert"
)

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
}
