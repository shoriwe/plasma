package simplification

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast"
	"reflect"
)

func SimplifyNode(node ast.Node) ast.Node {
	switch n := node.(type) {
	case ast.IExpression:
		return SimplifyExpression(n)
	case ast.IStatement:
		return SimplifyStatement(n)
	default:
		panic(fmt.Sprintf("unknown node type %s", reflect.TypeOf(node).String()))
	}
}

/*
SimplifyProgram pass does:
- Merges BEGIN and END to the program body
- Simplifies binary expressions
- Simplifies unary expression
- Transforms multiple result returns to tuple returns
- Transforms multiple result yields to tuple yields
- Transform unless statements to if statements
- Transforms unless one liners to if one liners
- Transforms for loops to while loops
- Transforms until loops to while loops
- Transforms generators definitions to classes
- Transforms interfaces to classes
*/
func SimplifyProgram(program *ast.Program) []ast.Node {
	length := len(program.Body)
	if program.Begin != nil {
		length += len(program.Begin.Body)
	}
	if program.End != nil {
		length += len(program.End.Body)
	}
	result := make([]ast.Node, 0, length)
	if program.Begin != nil {
		for _, child := range program.Begin.Body {
			result = append(result, SimplifyNode(child))
		}
	}
	for _, child := range program.Body {
		result = append(result, SimplifyNode(child))
	}
	if program.End != nil {
		for _, child := range program.End.Body {
			result = append(result, SimplifyNode(child))
		}
	}
	return result
}
