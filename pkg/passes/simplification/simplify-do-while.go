package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simp *simplify) simplifyDoWhile(do *ast.DoWhileStatement) *ast2.DoWhile {
	var (
		body      = make([]ast2.Node, 0, len(do.Body))
		condition = simp.simplifyExpression(do.Condition)
	)
	for _, node := range do.Body {
		body = append(body, simp.simplifyNode(node))
	}
	return &ast2.DoWhile{
		Body:      body,
		Condition: condition,
	}
}
