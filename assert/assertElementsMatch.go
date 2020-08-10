package assert

import (
	"testing"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/stretchr/testify/assert"
)

// AssertEqual dummy struct
type AssertElementsMatch struct {
}

func init() {
	function.Register(&AssertElementsMatch{})
}

// Name of function
func (s *AssertElementsMatch) Name() string {
	return "assertElementsMatch"
}

// Sig - function signature
func (s *AssertElementsMatch) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeArray, data.TypeArray}, false
}

// Eval - function implementation
func (s *AssertElementsMatch) Eval(params ...interface{}) (interface{}, error) {
	listA := params[0]
	listB := params[1]

	t := testing.T{}
	return assert.ElementsMatch(&t, listA, listB), nil
}
