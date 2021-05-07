package vm

import (
	"fmt"
	"testing"
)

func testMustSuccess(t *testing.T, samples map[string][]interface{}) {
	for expectedResult, code := range samples {
		vm := NewPlasmaVM()
		SetupDefaultTypes(vm)
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
		resultObject := result.(IObject)
		toString, getError := resultObject.Get(ToString)
		if getError != nil {
			t.Error(getError.String())
			return
		}
		var stringResult IObject
		stringResult, executionError = toString.(*Function).Callable.Call(vm.masterSymbolTable, vm, resultObject)
		if executionError != nil {
			t.Error(executionError.String())
			return
		}
		if stringResult.(*String).Value != expectedResult {
			t.Errorf("Expecting: %s but received: %s", expectedResult, stringResult.(*String).Value)
		}
	}
}

var newOPSamples = map[string][]interface{}{
	"Hello": {
		NewStringOP, "Hello",
		GetOP, StringName,
		CallOP, 1,
		ReturnOP,
	},
	"True": { // "Hello".ToBool()
		NewStringOP, "Hello",
		CopyOP, 2,
		GetFromOP, ToBool,
		CallOP, 1,
		ReturnOP,
	},
}

func TestData(t *testing.T) {
	testMustSuccess(t, newOPSamples)
}
