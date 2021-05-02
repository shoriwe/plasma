package vm

import (
	"testing"
)

func test(t *testing.T, code []interface{}, expect string) bool {
	vm := NewPlasmaVM(nil)
	vm.Initialize(code)
	result, executionError := vm.Execute()
	if executionError != nil {
		t.Error(executionError)
		return false
	}
	s, conversionError := result.String()
	if conversionError != nil {
		t.Error(conversionError.String())
		return false
	}
	finalResult := s.RawString() == expect
	if !finalResult {
		t.Errorf("Recevied: %s but expecting: %s", s.RawString(), expect)
	}
	return finalResult
}

var masterSymTable = NewSymbolTable(nil)

var binaryOperations = map[string][]interface{}{
	"18":  {PushOP, NewFloat(masterSymTable, "15"), PushOP, NewInteger(masterSymTable, "3", 10), AddOP, ReturnOP},
	"-30": {PushOP, NewInteger(masterSymTable, "50", 10), PushOP, NewInteger(masterSymTable, "20", 10), SubOP, ReturnOP},
	"4":   {PushOP, NewInteger(masterSymTable, "3", 10), PushOP, NewFloat(masterSymTable, "1"), AddOP, ReturnOP},
	"0.5": {PushOP, NewInteger(masterSymTable, "20", 10), PushOP, NewInteger(masterSymTable, "10", 10), DivOP, ReturnOP},
	"0.2": {PushOP, NewInteger(masterSymTable, "5", 10), PushOP, NewInteger(masterSymTable, "1", 10), DivOP, ReturnOP},
}

func TestBinaryOperations(t *testing.T) {
	for expect, sample := range binaryOperations {
		if !test(t, sample, expect) {
			return
		}
	}
}
