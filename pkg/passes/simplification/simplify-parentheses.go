package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func simplifyParentheses(expression *ast.ParenthesesExpression) ast2.Expression {
	return simplifyExpression(expression.X)
}
