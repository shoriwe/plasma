package vm

import (
	"fmt"
	"testing"
)

func test(t *testing.T, code []interface{}) *String {
	vm := NewPlasmaVM(nil)
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

var masterSymTable = NewSymbolTable(nil)

var binaryOperations = [][]interface{}{
	{PushOP, NewFloat(masterSymTable, "1000.6"), PushOP, NewInteger(masterSymTable, "13455", 10), AddOP},
}

func TestBinaryOperations(t *testing.T) {
	for _, sample := range binaryOperations {
		output := test(t, sample)
		if output == nil {
			return
		}
		fmt.Println(output.RawString())
	}
}
