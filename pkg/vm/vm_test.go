package vm

import (
	"fmt"
	"github.com/shoriwe/gruby/pkg/errors"
	"github.com/shoriwe/gruby/pkg/vm/object"
	_ "github.com/shoriwe/gruby/pkg/vm/object"
	_ "github.com/shoriwe/gruby/pkg/vm/utils"
	"testing"
)

func testMustSuccess(t *testing.T, samples map[string][]interface{}) {
	for expectedResult, code := range samples {
		vm := NewPlasmaVM()
		object.SetupDefaultTypes(vm)
		vm.Initialize(code)
		result, executionError := vm.Execute()
		if executionError != nil {
			t.Error(executionError.String())
			return
		}
		if result == nil {
			t.Error(fmt.Sprintf("Expecting %s, received nil", expectedResult))
			return
		}
		resultObject := result.(object.IObject)
		toString, getError := resultObject.Get(object.ToString)
		if getError != nil {
			t.Error(getError.String())
			return
		}
		var toStringCall interface{}
		toStringCall, getError = toString.(*object.Function).Get(object.Call)
		if getError != nil {
			t.Error(getError.String())
			return
		}
		rawStringObject, transformationError := toStringCall.(func(...object.IObject) (object.IObject, *errors.Error))(resultObject)
		if transformationError != nil {
			t.Error(transformationError.String())
			return
		}
		if rawStringObject.(*object.String).Value != expectedResult {
			t.Error(fmt.Sprintf("Expecting %s, received %s", expectedResult, rawStringObject.(*object.String).Value))
			return
		}

	}
}

var newOPSamples = map[string][]interface{}{
	"Hello": []interface{}{PushOP, "Hello", PushOP, 1, PushOP, object.StringName, GetOP, NewOP, ReturnOP},
}

func TestData(t *testing.T) {
	testMustSuccess(t, newOPSamples)
}
