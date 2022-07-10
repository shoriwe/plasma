package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simp *simplify) simplifyIndex(index *ast.IndexExpression) *ast2.Index {
	return &ast2.Index{
		Source: simp.simplifyExpression(index.Source),
		Index:  simp.simplifyExpression(index.Index),
	}
}
