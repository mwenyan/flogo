package report

import (
	"fmt"
	"testing"

	"github.com/project-flogo/core/app"
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

	input := make(map[string]string)
	input["testing1"] = "OK"
	input["testing2"] = "Failed"
	input["testing3"] = "abcd is not equal to abcd"
	app.SetValue("funitReport", input)
	_, err = act.Eval(tc)
	if err != nil {
		fmt.Errorf("has error: %v", err)
	}

	//check output
	output := tc.GetOutput("output")
	fmt.Printf("output=%v\n", output)

}
