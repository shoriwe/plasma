package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func simplifySelector(selector *ast.SelectorExpression) *ast2.Selector {
	return &ast2.Selector{
		X:          simplifyExpression(selector.X),
		Identifier: simplifyIdentifier(selector.Identifier),
	}
}
