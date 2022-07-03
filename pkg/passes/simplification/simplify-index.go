package simplification

import "github.com/shoriwe/gplasma/pkg/ast"

func simplifyIndex(index *ast.IndexExpression) *ast.IndexExpression {
	return &ast.IndexExpression{
		Source: SimplifyExpression(index.Source),
		Index:  SimplifyExpression(index.Index),
	}
}
