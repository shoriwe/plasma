package vm

import (
	"fmt"
	"github.com/shoriwe/gruby/pkg/vm/object"
	"testing"
)

func test(t *testing.T, code []interface{}) string {
	vm := NewPlasmaVM(nil)
	result, executionError := vm.Execute(code)
	if executionError != nil {
		t.Error(executionError)
		return "ERROR"
	}
	s, conversionError := result.RawString()
	if conversionError != nil {
		t.Error(conversionError.String())
	}
	return s
}

var binaryOperations = [][]interface{}{
	{PushOP, object.NewInteger("1000", 10), PushOP, object.NewInteger("13455", 10), AddOP},
}

func TestBinaryOperations(t *testing.T) {
	for _, sample := range binaryOperations {
		output := test(t, sample)
		if output == "ERROR" {
			return
		}
		fmt.Println(output)
	}
}
