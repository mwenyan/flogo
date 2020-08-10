package assert

import (
	"testing"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/stretchr/testify/assert"
)

// AssertEqual dummy struct
type AssertJSONEq struct {
}

func init() {
	function.Register(&AssertJSONEq{})
}

// Name of function
func (s *AssertJSONEq) Name() string {
	return "assertJSONEq"
}

// Sig - function signature
func (s *AssertJSONEq) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeString}, false
}

// Eval - function implementation
func (s *AssertJSONEq) Eval(params ...interface{}) (interface{}, error) {
	a := params[0].(string)
	b := params[1].(string)

	t := testing.T{}
	return assert.JSONEq(&t, a, b), nil
}
