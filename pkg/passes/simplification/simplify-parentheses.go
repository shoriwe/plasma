package simplification

import (
	"github.com/shoriwe/plasma/pkg/ast"
	"github.com/shoriwe/plasma/pkg/ast2"
)

func (simplify *simplifyPass) Parentheses(expression *ast.ParenthesesExpression) ast2.Expression {
	return simplify.Expression(expression.X)
}
