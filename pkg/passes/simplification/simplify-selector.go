package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simp *simplify) simplifySelector(selector *ast.SelectorExpression) *ast2.Selector {
	return &ast2.Selector{
		X:          simp.simplifyExpression(selector.X),
		Identifier: simp.simplifyIdentifier(selector.Identifier),
	}
}
