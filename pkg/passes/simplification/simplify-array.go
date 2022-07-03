package simplification

import "github.com/shoriwe/gplasma/pkg/ast"

func simplifyArray(array *ast.ArrayExpression) *ast.ArrayExpression {
	newContents := make([]ast.IExpression, 0, len(array.Values))
	for _, value := range array.Values {
		newContents = append(newContents, SimplifyExpression(value))
	}
	return &ast.ArrayExpression{
		Values: newContents,
	}
}
