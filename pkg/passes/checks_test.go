package passes

import (
	_ "embed"
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/lexer"
	"github.com/shoriwe/gplasma/pkg/parser"
	"github.com/shoriwe/gplasma/pkg/reader"
	"testing"
)

var (
	//go:embed valid-script.pm
	validScript string
	//go:embed invalid-return.pm
	invalidFunctionScript string
	//go:embed invalid-yield.pm
	invalidYieldScript string
	//go:embed invalid-redo.pm
	invalidRedoScript string
	//go:embed invalid-continue.pm
	invalidContinueScript string
	//go:embed invalid-break.pm
	invalidBreakScript string
)

func executeScript(script string) (*Check, error) {
	checkPass := NewCheckPass()
	codeReader := reader.NewStringReader(script)
	l := lexer.NewLexer(codeReader)
	p := parser.NewParser(l)
	program, parseError := p.Parse()
	if parseError != nil {
		return nil, parseError
	}
	ast.Walk(checkPass, program)
	return checkPass, nil
}

func TestValidScript(t *testing.T) {
	checkPass, passError := executeScript(validScript)
	if passError != nil {
		t.Fatal(passError)
	}
	if checkPass.CountInvalidFunctionNodes() > 0 {
		t.Fatal("Invalid returns found")
	}
	if checkPass.CountInvalidGeneratorNodes() > 0 {
		t.Fatal("Invalid Yields found")
	}
	if checkPass.CountInvalidLoopNodes() > 0 {
		t.Fatal("Invalid break/redo/continue found")
	}
}

func TestInvalidReturn(t *testing.T) {
	checkPass, passError := executeScript(invalidFunctionScript)
	if passError != nil {
		t.Fatal(passError)
	}
	const invalidReturns = 3
	if found := checkPass.CountInvalidFunctionNodes(); found != invalidReturns {
		t.Fatalf("Invalid returns bypassed the check pass, expecting %d but found %d", invalidReturns, found)
	}
}

func TestInvalidYield(t *testing.T) {
	checkPass, passError := executeScript(invalidYieldScript)
	if passError != nil {
		t.Fatal(passError)
	}
	const invalidReturns = 4
	if found := checkPass.CountInvalidGeneratorNodes(); found != invalidReturns {
		t.Fatalf("Invalid returns bypassed the check pass, expecting %d but found %d", invalidReturns, found)
	}
}

func TestInvalidRedo(t *testing.T) {
	checkPass, passError := executeScript(invalidRedoScript)
	if passError != nil {
		t.Fatal(passError)
	}
	const invalidCount = 5
	if found := checkPass.CountInvalidLoopNodes(); found != invalidCount {
		t.Fatalf("Invalid returns bypassed the check pass, expecting %d but found %d", invalidCount, found)
	}
}

func TestInvalidContinue(t *testing.T) {
	checkPass, passError := executeScript(invalidContinueScript)
	if passError != nil {
		t.Fatal(passError)
	}
	const invalidCount = 5
	if found := checkPass.CountInvalidLoopNodes(); found != invalidCount {
		t.Fatalf("Invalid returns bypassed the check pass, expecting %d but found %d", invalidCount, found)
	}
}

func TestInvalidBreak(t *testing.T) {
	checkPass, passError := executeScript(invalidBreakScript)
	if passError != nil {
		t.Fatal(passError)
	}
	const invalidCount = 5
	if found := checkPass.CountInvalidLoopNodes(); found != invalidCount {
		t.Fatalf("Invalid returns bypassed the check pass, expecting %d but found %d", invalidCount, found)
	}
}
