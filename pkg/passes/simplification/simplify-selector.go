package simplification

import (
	"github.com/shoriwe/plasma/pkg/ast"
	"github.com/shoriwe/plasma/pkg/ast2"
)

func (simplify *simplifyPass) Selector(selector *ast.SelectorExpression) *ast2.Selector {
	return &ast2.Selector{
		X:          simplify.Expression(selector.X),
		Identifier: simplify.Identifier(selector.Identifier),
	}
}
