package compiler

import (
	"fmt"
	"github.com/shoriwe/plasma/pkg/ast"
	"github.com/shoriwe/plasma/pkg/bytecode/assembler"
	"github.com/shoriwe/plasma/pkg/lexer"
	"github.com/shoriwe/plasma/pkg/parser"
	"github.com/shoriwe/plasma/pkg/passes/checks"
	"github.com/shoriwe/plasma/pkg/passes/simplification"
	transformations_1 "github.com/shoriwe/plasma/pkg/passes/transformations-1"
	"github.com/shoriwe/plasma/pkg/reader"
)

func Compile(scriptCode string) ([]byte, error) {
	l := lexer.NewLexer(reader.NewStringReader(scriptCode))
	p := parser.NewParser(l)
	programAst1, parseError := p.Parse()
	if parseError != nil {
		return nil, parseError
	}
	checkPass := checks.NewCheckPass()
	ast.Walk(checkPass, programAst1)
	if checkPass.CountInvalidLoopNodes() > 0 {
		return nil, fmt.Errorf("invalid loop nodes found")
	}
	if checkPass.CountInvalidFunctionNodes() > 0 {
		return nil, fmt.Errorf("invalid function nodes found")
	}
	if checkPass.CountInvalidGeneratorNodes() > 0 {
		return nil, fmt.Errorf("invalid generator nodes found")
	}
	programAst2, simplifyError := simplification.Simplify(programAst1)
	if simplifyError != nil {
		return nil, simplifyError
	}
	programAst3, transformError := transformations_1.Transform(programAst2)
	if transformError != nil {
		return nil, transformError
	}
	return assembler.Assemble(programAst3)
}
