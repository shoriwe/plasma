package simplification

import "github.com/shoriwe/gplasma/pkg/ast"

func simplifyTuple(tuple *ast.TupleExpression) *ast.TupleExpression {
	newContents := make([]ast.IExpression, 0, len(tuple.Values))
	for _, value := range tuple.Values {
		newContents = append(newContents, SimplifyExpression(value))
	}
	return &ast.TupleExpression{
		Values: newContents,
	}
}
