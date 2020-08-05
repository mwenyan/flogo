package mock

import (
	"encoding/json"
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
	mockValue := `{"test":"this is testing"}`
	obj := make(map[string]interface{})
	err := json.Unmarshal([]byte(mockValue), &obj)
	if err != nil {
		fmt.Errorf("has error: %v", err)
	}

	iCtx := test.NewActivityInitContext(Settings{MockValue: obj}, mf)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())

	_, err = act.Eval(tc)
	if err != nil {
		fmt.Errorf("has error: %v", err)
	}

	//check output
	output := tc.GetOutput("output")
	fmt.Printf("output=%v\n", output)

}