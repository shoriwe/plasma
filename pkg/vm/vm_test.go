package vm

import (
	"github.com/shoriwe/gruby/pkg/vm/runtime"
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

var masterSymTable = runtime.NewSymbolTable(nil)

var binaryOperations = map[string][]interface{}{
	"18":    {runtime.PushOP, runtime.NewFloat(masterSymTable, "15"), runtime.PushOP, runtime.NewInteger(masterSymTable, "3", 10), runtime.AddOP, runtime.ReturnOP},
	"-30":   {runtime.PushOP, runtime.NewInteger(masterSymTable, "50", 10), runtime.PushOP, runtime.NewInteger(masterSymTable, "20", 10), runtime.SubOP, runtime.ReturnOP},
	"-29.5": {runtime.PushOP, runtime.NewInteger(masterSymTable, "50", 10), runtime.PushOP, runtime.NewFloat(masterSymTable, "20.5"), runtime.SubOP, runtime.ReturnOP},
	"4":     {runtime.PushOP, runtime.NewInteger(masterSymTable, "3", 10), runtime.PushOP, runtime.NewFloat(masterSymTable, "1"), runtime.AddOP, runtime.ReturnOP},
	"0.5":   {runtime.PushOP, runtime.NewInteger(masterSymTable, "20", 10), runtime.PushOP, runtime.NewInteger(masterSymTable, "10", 10), runtime.DivOP, runtime.ReturnOP},
	"0.2":   {runtime.PushOP, runtime.NewInteger(masterSymTable, "5", 10), runtime.PushOP, runtime.NewInteger(masterSymTable, "1", 10), runtime.DivOP, runtime.ReturnOP},
	"5.6":   {runtime.PushOP, runtime.NewFloat(masterSymTable, "1"), runtime.PushOP, runtime.NewFloat(masterSymTable, "5.6"), runtime.MulOP, runtime.ReturnOP},
	"20":    {runtime.PushOP, runtime.NewFloat(masterSymTable, "2"), runtime.PushOP, runtime.NewInteger(masterSymTable, "10", 10), runtime.MulOP, runtime.ReturnOP},
	"0.6":   {runtime.PushOP, runtime.NewFloat(masterSymTable, "2"), runtime.PushOP, runtime.NewFloat(masterSymTable, "0.3"), runtime.MulOP, runtime.ReturnOP},
}

func TestBinaryOperations(t *testing.T) {
	for expect, sample := range binaryOperations {
		if !test(t, sample, expect) {
			return
		}
	}
}
