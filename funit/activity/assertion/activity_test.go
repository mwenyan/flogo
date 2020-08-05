package assert

import (
	"fmt"
	"testing"

	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/resolve"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateDoc(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	mf := mapper.NewFactory(resolve.GetBasicResolver())
	iCtx := test.NewActivityInitContext(Settings{}, mf)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())
	//	inputmp := make(map[string]interface{})
	//	inputmp["Value"] = `{"observers":[],"issuer":"Alice","amount":"10000","currency":"USD","owner":"Alice"}`
	//	v, err := json.Marshal(inputmp)
	//	complexObject := &data.ComplexObject{}
	//	err = json.Unmarshal([]byte(v), complexObject)

	input := make(map[string]interface{})
	input["assertion"] = false
	input["msg_var0"] = "abc"
	input["msg_var1"] = "abcd"
	err = tc.SetInputObject(&Input{Data: input, Name: "verify string value", ErrorMessage: "expected value %v is not equal to actual value %v"})
	if err != nil {
		fmt.Errorf("set input error: %v", err)
	}
	_, err = act.Eval(tc)
	if err != nil {
		fmt.Errorf("has error: %v", err)
	}

	//check output
	output := tc.GetOutput("output")
	fmt.Printf("output=%v\n", output)

}
