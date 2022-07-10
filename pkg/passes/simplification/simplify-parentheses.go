package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simp *simplify) simplifyParentheses(expression *ast.ParenthesesExpression) ast2.Expression {
	return simp.simplifyExpression(expression.X)
}
