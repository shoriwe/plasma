package checks

import (
	_ "embed"
	"github.com/shoriwe/plasma/pkg/ast"
	"github.com/shoriwe/plasma/pkg/lexer"
	"github.com/shoriwe/plasma/pkg/parser"
	"github.com/shoriwe/plasma/pkg/reader"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	//go:embed valid-script.pm
	validScript string
	//go:embed invalid-return.pm
	invalidFunctionScript string
	//go:embed invalid-yield.pm
	invalidYieldScript string
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
	assert.Nil(t, passError)
	assert.Equal(t, 0, checkPass.CountInvalidFunctionNodes(), "Invalid returns found")
	assert.Equal(t, 0, checkPass.CountInvalidGeneratorNodes(), "Invalid Yields found")
	assert.Equal(t, 0, checkPass.CountInvalidLoopNodes(), "Invalid break/redo/continue found")
}

func TestInvalidReturn(t *testing.T) {
	checkPass, passError := executeScript(invalidFunctionScript)
	assert.Nil(t, passError)
	assert.Equal(t, 4, checkPass.CountInvalidFunctionNodes())
}

func TestInvalidYield(t *testing.T) {
	checkPass, passError := executeScript(invalidYieldScript)
	assert.Nil(t, passError)
	assert.Equal(t, 4, checkPass.CountInvalidGeneratorNodes())
}

func TestInvalidContinue(t *testing.T) {
	checkPass, passError := executeScript(invalidContinueScript)
	assert.Nil(t, passError)
	assert.Equal(t, 5, checkPass.CountInvalidLoopNodes())
}

func TestInvalidBreak(t *testing.T) {
	checkPass, passError := executeScript(invalidBreakScript)
	assert.Nil(t, passError)
	assert.Equal(t, 5, checkPass.CountInvalidLoopNodes())
}
