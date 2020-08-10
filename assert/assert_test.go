package assert

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertEqual(t *testing.T) {
	aeq := &AssertEqual{}
	r, err := aeq.Eval(1, 1)
	assert.True(t, r.(bool))
	assert.Nil(t, err)

	m1 := make(map[string]interface{})
	m1["a"] = "1"
	m1["b"] = 2
	m1["c"] = true

	m2 := make(map[string]interface{})
	m2["a"] = "1"
	m2["b"] = 2
	m2["c"] = false
	assert.Equal(t, m1, m2)
}

func TestColor(t *testing.T) {
	//#ffa500
	c, err := getColorName("#ffa500")
	assert.Nil(t, err)

	fmt.Println(c)

	c, _ = getColorName("#ff0000")
	fmt.Println(c)

}
