package assert

import (
	"testing"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/stretchr/testify/assert"
)

// AssertEqual dummy struct
type AssertSubset struct {
}

func init() {
	function.Register(&AssertSubset{})
}

// Name of function
func (s *AssertSubset) Name() string {
	return "assertSubset"
}

// Sig - function signature
func (s *AssertSubset) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeArray, data.TypeArray}, false
}

// Eval - function implementation
func (s *AssertSubset) Eval(params ...interface{}) (interface{}, error) {
	listA := params[0]
	listB := params[1]

	t := testing.T{}
	return assert.Subset(&t, listA, listB), nil
}
