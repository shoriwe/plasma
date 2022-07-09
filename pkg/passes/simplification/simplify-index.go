package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func simplifyIndex(index *ast.IndexExpression) *ast2.Index {
	return &ast2.Index{
		Source: simplifyExpression(index.Source),
		Index:  simplifyExpression(index.Index),
	}
}
