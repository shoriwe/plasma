package vm

import (
	"fmt"
	"github.com/shoriwe/gruby/pkg/vm/object"
	"testing"
)

func test(t *testing.T, code []interface{}) *object.String {
	vm := NewPlasmaVM()
	result, executionError := vm.Execute(code)
	if executionError != nil {
		t.Error(executionError)
		return nil
	}
	s, conversionError := result.String()
	if conversionError != nil {
		t.Error(conversionError.String())
	}
	return s
}

var binaryOperations = [][]interface{}{
	{nil, nil},
}

func TestBinaryOperations(t *testing.T) {
	for _, sample := range binaryOperations {
		output := test(t, sample)
		if output == nil {
			return
		}
		fmt.Println(output)
	}
}
