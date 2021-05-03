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
	s, conversionError := result.RawString()
	if conversionError != nil {
		t.Error(conversionError.String())
		return false
	}
	finalResult := s == expect
	if !finalResult {
		t.Errorf("Recevied: %s but expecting: %s", s, expect)
	}
	return finalResult
}

var masterSymTable = runtime.NewSymbolTable(nil)

var divOperations = map[string][]interface{}{
	"5.000000": {runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "100"), runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "500"), runtime.DivOP, runtime.ReturnOP},
	"7.000000": {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "3", 10), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "21", 10), runtime.DivOP, runtime.ReturnOP},
	"0.000000": {runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "1"), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "0", 10), runtime.DivOP, runtime.ReturnOP},
	"2.500000": {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "1", 10), runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "2.5"), runtime.DivOP, runtime.ReturnOP},
}

var mulOperations = map[string][]interface{}{
	"5.000000":  {runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "0.1"), runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "50"), runtime.MulOP, runtime.ReturnOP},
	"0":         {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "3", 10), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "0", 10), runtime.MulOP, runtime.ReturnOP},
	"-3.000000": {runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "-1"), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "3", 10), runtime.MulOP, runtime.ReturnOP},
	"9.000000":  {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "3", 10), runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "3"), runtime.MulOP, runtime.ReturnOP},
}

var powOperations = map[string][]interface{}{
	"4.000000":  {runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "2"), runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "2"), runtime.PowOP, runtime.ReturnOP},
	"9":         {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "2", 10), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "3", 10), runtime.PowOP, runtime.ReturnOP},
	"2.000000":  {runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "0.5"), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "4", 10), runtime.PowOP, runtime.ReturnOP},
	"27.000000": {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "3", 10), runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "3"), runtime.PowOP, runtime.ReturnOP},
}
var subOperations = map[string][]interface{}{
	"-11.000000": {runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "15"), runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "4"), runtime.SubOP, runtime.ReturnOP},
	"-10":        {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "15", 10), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "5", 10), runtime.SubOP, runtime.ReturnOP},
	"-30.000000": {runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "50"), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "20", 10), runtime.SubOP, runtime.ReturnOP},
	"-29.500000": {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "50", 10), runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "20.5"), runtime.SubOP, runtime.ReturnOP},
}
var addOperations = map[string][]interface{}{
	"19.000000": {runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "15"), runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "4"), runtime.AddOP, runtime.ReturnOP},
	"20":        {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "15", 10), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "5", 10), runtime.AddOP, runtime.ReturnOP},
	"18.000000": {runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "15"), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "3", 10), runtime.AddOP, runtime.ReturnOP},
	"4.000000":  {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "3", 10), runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "1"), runtime.AddOP, runtime.ReturnOP},
}
var modOperations = map[string][]interface{}{
	"0.000000": {runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "5"), runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "15"), runtime.ModOP, runtime.ReturnOP},
	"1":        {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "9", 10), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "10", 10), runtime.ModOP, runtime.ReturnOP},
	"5.000000": {runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "6"), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "5", 10), runtime.ModOP, runtime.ReturnOP},
	"5.900000": {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "6", 10), runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "5.9"), runtime.ModOP, runtime.ReturnOP},
}
var floorDivOperations = map[string][]interface{}{
	"0":  {runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "15"), runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "5"), runtime.FloorDivOP, runtime.ReturnOP},
	"1":  {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "9", 10), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "10", 10), runtime.FloorDivOP, runtime.ReturnOP},
	"15": {runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "5"), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "75", 10), runtime.FloorDivOP, runtime.ReturnOP},
	"5":  {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "1", 10), runtime.PushOP, runtime.NewFloatFromString(masterSymTable, "5"), runtime.FloorDivOP, runtime.ReturnOP},
}
var bitwiseLeftOperations = map[string][]interface{}{
	"163840": {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "15", 10), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "5", 10), runtime.BitwiseLeftOP, runtime.ReturnOP},
}
var bitwiseRightOperations = map[string][]interface{}{
	"0": {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "15", 10), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "5", 10), runtime.BitwiseRightOP, runtime.ReturnOP},
}
var bitwiseAndOperations = map[string][]interface{}{
	"5": {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "15", 10), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "5", 10), runtime.BitwiseAndOP, runtime.ReturnOP},
}
var bitwiseOrOperations = map[string][]interface{}{
	"0": {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "15346", 10), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "5", 10), runtime.BitwiseOrOP, runtime.ReturnOP},
}
var bitwiseXorOperations = map[string][]interface{}{
	"5": {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "4", 10), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "1", 10), runtime.BitwiseXorOP, runtime.ReturnOP},
}
var andIntOnlyOperations = map[string][]interface{}{
	"True":  {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "4", 10), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "1", 10), runtime.AndOP, runtime.ReturnOP},
	"False": {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "4", 10), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "0", 10), runtime.AndOP, runtime.ReturnOP},
}
var andBoolOnlyOperations = map[string][]interface{}{
	"True":  {runtime.PushOP, runtime.NewTrue(masterSymTable), runtime.PushOP, runtime.NewTrue(masterSymTable), runtime.AndOP, runtime.ReturnOP},
	"False": {runtime.PushOP, runtime.NewFalse(masterSymTable), runtime.PushOP, runtime.NewTrue(masterSymTable), runtime.AndOP, runtime.ReturnOP},
}
var orIntOnlyOperations = map[string][]interface{}{
	"True":  {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "4", 10), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "1", 10), runtime.OrOP, runtime.ReturnOP},
	"False": {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "0", 10), runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "0", 10), runtime.OrOP, runtime.ReturnOP},
}
var orBoolOnlyOperations = map[string][]interface{}{
	"True":  {runtime.PushOP, runtime.NewTrue(masterSymTable), runtime.PushOP, runtime.NewTrue(masterSymTable), runtime.OrOP, runtime.ReturnOP},
	"False": {runtime.PushOP, runtime.NewFalse(masterSymTable), runtime.PushOP, runtime.NewFalse(masterSymTable), runtime.OrOP, runtime.ReturnOP},
}
var negateIntOnlyOperations = map[string][]interface{}{
	"False":  {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "1", 10), runtime.NegateOP, runtime.ReturnOP},
	"True": {runtime.PushOP, runtime.NewIntegerFromString(masterSymTable, "0", 10), runtime.NegateOP, runtime.ReturnOP},
}
var negateBoolOnlyOperations = map[string][]interface{}{
	"False":  {runtime.PushOP, runtime.NewTrue(masterSymTable), runtime.NegateOP, runtime.ReturnOP},
	"True": {runtime.PushOP, runtime.NewFalse(masterSymTable), runtime.NegateOP, runtime.ReturnOP},
}

func TestDivOperations(t *testing.T) {
	for expect, sample := range divOperations {
		if !test(t, sample, expect) {
			return
		}
	}
}

func TestMulOperations(t *testing.T) {
	for expect, sample := range mulOperations {
		if !test(t, sample, expect) {
			return
		}
	}
}

func TestAddOperations(t *testing.T) {
	for expect, sample := range addOperations {
		if !test(t, sample, expect) {
			return
		}
	}
}

func TestSubOperations(t *testing.T) {
	for expect, sample := range subOperations {
		if !test(t, sample, expect) {
			return
		}
	}
}

func TestPowOperations(t *testing.T) {
	for expect, sample := range powOperations {
		if !test(t, sample, expect) {
			return
		}
	}
}

func TestModOperations(t *testing.T) {
	for expect, sample := range modOperations {
		if !test(t, sample, expect) {
			return
		}
	}
}

func TestFloorDivOperations(t *testing.T) {
	for expect, sample := range floorDivOperations {
		if !test(t, sample, expect) {
			return
		}
	}
}

func TestBitwiseLeftOperations(t *testing.T) {
	for expect, sample := range bitwiseLeftOperations {
		if !test(t, sample, expect) {
			return
		}
	}
}

func TestBitwiseRightOperations(t *testing.T) {
	for expect, sample := range bitwiseRightOperations {
		if !test(t, sample, expect) {
			return
		}
	}
}

func TestBitwiseAndOperations(t *testing.T) {
	for expect, sample := range bitwiseAndOperations {
		if !test(t, sample, expect) {
			return
		}
	}
}

func TestBitwiseOrOperations(t *testing.T) {
	for expect, sample := range bitwiseOrOperations {
		if !test(t, sample, expect) {
			return
		}
	}
}

func TestBitwiseXorOperations(t *testing.T) {
	for expect, sample := range bitwiseXorOperations {
		if !test(t, sample, expect) {
			return
		}
	}
}

func TestAndIntOnlyOperations(t *testing.T) {
	for expect, sample := range andIntOnlyOperations {
		if !test(t, sample, expect) {
			return
		}
	}
}

func TestAndBoolOnlyOperations(t *testing.T) {
	for expect, sample := range andBoolOnlyOperations {
		if !test(t, sample, expect) {
			return
		}
	}
}

func TestOrIntOnlyOperations(t *testing.T) {
	for expect, sample := range orIntOnlyOperations {
		if !test(t, sample, expect) {
			return
		}
	}
}

func TestOrBoolOnlyOperations(t *testing.T) {
	for expect, sample := range orBoolOnlyOperations {
		if !test(t, sample, expect) {
			return
		}
	}
}

func TestNegateIntOnlyOperations(t *testing.T) {
	for expect, sample := range negateIntOnlyOperations {
		if !test(t, sample, expect) {
			return
		}
	}
}

func TestNegateBoolOnlyOperations(t *testing.T) {
	for expect, sample := range negateBoolOnlyOperations {
		if !test(t, sample, expect) {
			return
		}
	}
}
